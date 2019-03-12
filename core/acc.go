package core

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"tribc/lib"
)

/*
  author: Guo Guisheng
  本文件主要完成了对账号结构的定义，并实现了针对账号交互的关键操作
 */

const (
	version            = "1.1"
)

//密钥对结构
type GKey struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

//账户结构（为适配隐私地址技术，需要两个公私钥对，若不实现隐私地址技术，可不用）
type Account struct {
	GkeyA *GKey
	GkeyB *GKey
}


//根据输入随机字符串，生成一个密钥对
func MakeNewKey(randKey string) (*GKey, error) {
	var err error
	var gkey GKey

	private, err := ecdsa.GenerateKey(curve, strings.NewReader(randKey))
	if err != nil {
		log.Panic(err)
	}
	gkey = GKey{private, private.PublicKey}
	return &gkey, nil
}


//创建账号，根据当前的设计业务设计，所有新建的account都需要通过调用此函数进行生成，否则此账号对象在签名验证过程中会有错误（本质上通过调用此函数，实现GkeyA,GkeyB一致）
func CreateAccount(randomstr string) *Account{
	gkeyA, err := MakeNewKey(randomstr)
	gkeyB, err := MakeNewKey(randomstr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	privKeyA := gkeyA.GetPrivKey()
	privKeyB := gkeyB.GetPrivKey()
	fmt.Println("A privateKey is :", lib.ByteToString(privKeyA))
	//fmt.Println("A privateKey is :",  hex.EncodeToString(privKeyA)) lib.ByteToString hex.EncodeToString 为等效功能函数
	fmt.Println("B privateKey is :", lib.ByteToString(privKeyB))
	pubKeyA := gkeyA.GetPubKey()
	pubKeyB := gkeyB.GetPubKey()
	fmt.Println("A publickKey is :", lib.ByteToString(pubKeyA))
	fmt.Println("B publickKey is :", lib.ByteToString(pubKeyB))

	acc := Account{gkeyA,gkeyB}
	return &acc
}

//根据私钥byte获得gkey
func Priv2gkey(priv_s []byte) *GKey  {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(priv_s)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(priv_s)
	gkey := GKey{priv, priv.PublicKey}
	return &gkey
}

//根据公钥字节串获得一个公钥对象
func Pub2pubKey(pub_b []byte) *ecdsa.PublicKey{
	//pub_b,_:=hex.DecodeString(s)
	if len(pub_b)==64{
		pubkey:=ecdsa.PublicKey{curve,new(big.Int).SetBytes(pub_b[:32]),new(big.Int).SetBytes(pub_b[32:])}
		return &pubkey
	}
	return nil
}


//根据公钥对，返回私钥byte
func (k GKey) GetPrivKey() []byte {
	d := k.PrivateKey.D.Bytes()
	b := make([]byte, 0, 32)
	priKey := lib.PaddedAppend(32, b, d) // []bytes type
	// s := byteToString(priKey)
	return priKey
}

//根据公钥对，返回公钥byte
func (k GKey) GetPubKey() []byte {
	pubKey := append(PubkeyPad(k.PublicKey.X.Bytes()), PubkeyPad(k.PublicKey.Y.Bytes())...)//k.PrivateKey.Y.Bytes()...) // []bytes type
	// s := byteToString(pubKey)
	return pubKey
}

//根据隐私地址账户获得用户操作地址 （若不适用隐私地址的双层公私钥账户，可自定义修改实现
func GetAddress(pub_bytes []byte) (address string) {
	/*
	// SHA256 HASH
	//fmt.Println("1 - Perform SHA-256 hashing on the public key")
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(pub_bytes)
	pub_hash_1 := sha256_h.Sum(nil) // 对公钥进行hash256运算
	//fmt.Println(lib.ByteToString(pub_hash_1))
	//fmt.Println("================")

	// RIPEMD-160 HASH
	//fmt.Println("2 - Perform RIPEMD-160 hashing on the result of SHA-256")
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil) // 对公钥hash进行ripemd160运算
	*/
	//fmt.Println(lib.ByteToString(pub_hash_2))
	//fmt.Println("================")
	/* Convert hash bytes to base58 chech encoded sequence */
	//address = lib.B58checkencode(0x00, pub_hash_2)

	return lib.B58encode(pub_bytes)
}

//根据地址返回公钥对象
func GetPubk4Addr(addr string) *ecdsa.PublicKey{
	pub_b:=lib.Base58Decode(addr)
	pub:=Pub2pubKey(pub_b)
	return pub
}

/*
对数据进行签名，其中hash为需要签名数据的哈希
返回加密结果，结果为数字证书r、s的序列化后拼接，然后用hex转换为string
*/
func Sign(priv *ecdsa.PrivateKey,hash []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return hex.EncodeToString(b.Bytes()), nil
}



/*
签名数据分解恢复R,S
*/
func getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		err = errors.New("decrypt error," + err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

/*
校验文本内容是否与签名一致
使用公钥校验签名和文本内容
*/
func Verify(hash []byte, signature string, pubKey *ecdsa.PublicKey) (bool, error) {
	rint, sint, err := getSign(signature)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(pubKey, hash, &rint, &sint)
	return result, nil
}

type AccAes struct{
	Privhash string
	Privaes string
}
type Accfile struct {
	A AccAes
	B AccAes
	V string
	Name string
}

//保存账号到文件，采用加密方式保存{a:{privhash:hash(acc_a),privaes:aes(acc_a)},b:{privhash:hash(acc_b),privaes:aes(acc_b)}
func Save2file(acc *Account,path string,key []byte) bool{
	var af Accfile
	af.V=version
	priv_a := acc.GkeyA.GetPrivKey()
	priv_b := acc.GkeyB.GetPrivKey()
	md5key:=md5.Sum(priv_a)
	af.A.Privhash=hex.EncodeToString(md5key[:])
	xpass, err :=lib.AesEncrypt(priv_a,key)
	if err != nil {
		fmt.Println("[Error Save2file]",err)
		return false
	}
	af.A.Privaes=base64.StdEncoding.EncodeToString(xpass)

	md5key=md5.Sum(priv_b)
	af.B.Privhash=hex.EncodeToString(md5key[:])
	xpass, err =lib.AesEncrypt(priv_b,key)
	if err != nil {
		fmt.Println("[Error Save2file]",err)
		return false
	}
	af.B.Privaes=base64.StdEncoding.EncodeToString(xpass)
	af.Name=GetAddress(acc.GkeyA.GetPubKey())

	data,_:= json.Marshal(af)
	if ioutil.WriteFile(path,data,0644)==nil{
		fmt.Println("[Info Save2file]","写入账号文件成功",path)
		return true
	}
	return false
}

//导入账号文件
func Load4file(path string,key []byte) *Account{
	var af Accfile
	var acc Account
	data,_ := ioutil.ReadFile(path)
	json.Unmarshal(data,&af)
	acc.GkeyA,_ = entry2gkey(&af.A,key)
	acc.GkeyB,_ = entry2gkey(&af.B,key)
	if acc.GkeyA==nil || acc.GkeyB==nil{
		fmt.Println("[Error load4file] 文件或密码错误，导入账号文件失败")
		return nil
	}
	fmt.Println("[Info Load4file]","导入账号文件成功",path)
	return &acc
}

//加密账号
func entry2gkey(aa * AccAes,key []byte) (*GKey,error) {
	bytesPass, err := base64.StdEncoding.DecodeString(aa.Privaes)
	priv, err := lib.AesDecrypt(bytesPass,key)
	if priv==nil || err != nil {
		fmt.Println("[Error entry2gkey]","密码错误，解密失败")
		return nil,errors.New("解密失败")
	}

	md5key:=md5.Sum(priv)
	if hex.EncodeToString(md5key[:])== aa.Privhash{
		return Priv2gkey(priv),nil
	}
	return nil,errors.New("密码错误")
}


//根据公钥返回地址，用于在区块链中检验当前签名是否是当前地址的签名（签名检查分两步：签名有效性检查，当前签名用户地址和当前资产地址检查）
func Pubkey2Addr(puk_s string)string{
	pub_b,_:= hex.DecodeString(puk_s)
	result:=GetAddress(pub_b)
	return result
}

//根据公钥返回地址(仅用于隐私交易下隐藏地址交易转化），用于在区块链中检验当前签名是否是当前地址的签名
func Shield_Pubkey2Addr(s_puk_s string)string{
	pub_b,_:= hex.DecodeString(s_puk_s)
	result:=hex.EncodeToString(pub_b[:32])
	return result
}