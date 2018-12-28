package lib

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/ripemd160"
)

// 把字节数组转换为字符串
func ByteToString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}
	return s
}


func PaddedAppend(size uint, dst, src []byte) []byte {
	/*
		把src数组转换成指定长度的数组，长度不够则添加0

			:param size: 要返回的数组长度
			:param dst: byte类型的切片，需要返回的切片
			:param src: 原byte数组
	*/
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// 把pubkey进行二次哈希运算
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey) // 第一次哈希运算
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil) // 第二次哈希运算
	return publicRIPEMD160
}

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// b58encode encodea a byte slice b into a base-58 encoded string.
func B58encode(b []byte) (s string) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	x := new(big.Int).SetBytes(b)
	// Initialize
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x /58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(b58Alphabet[r.Int64()]) + s
	}
	return s
}

func Base58Decode(s string) []byte {
	input:=[]byte(s)
	result := big.NewInt(0)
	zeroBytes := 0

	for _, b := range input {
		if b != b58Alphabet[0] {
			break
		}

		zeroBytes++
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(int64(len(b58Alphabet))))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}


// b58chechencode encodes version ver and byte slice b into a base-58 chech encoded string.
func B58checkencode(ver uint8, b []byte) (s string) {
	/* Prepend version */
	fmt.Println("3 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)")
	bcpy := append([]byte{ver}, b...)
	fmt.Println(ByteToString(bcpy))
	fmt.Println("================")

	/* Create a new SHA256 context */
	sha256H := sha256.New()

	/* SHA256 HASH #1 */
	fmt.Println("4 - Perform SHA-256 hash on the extended PIPEMD-160 result")
	sha256H.Reset()
	sha256H.Write(bcpy)
	hash1 := sha256H.Sum(nil)
	fmt.Println(ByteToString(hash1))
	fmt.Println("================")

	/* SHA256 HASH #2 */
	fmt.Println("5 - Perform SHA-256 hash on the result of the previous SHA-256 hash")
	sha256H.Reset()
	sha256H.Write(hash1)
	hash2 := sha256H.Sum(nil)
	fmt.Println(ByteToString(hash2))
	fmt.Println("================")

	/* Append first four bytes of hash */
	fmt.Println("6 - Take the first 4 bytes of the second SHA-256 hash. This is the address chechsum")
	fmt.Println(ByteToString(hash2[0:4]))
	fmt.Println("================")

	fmt.Println("7 - Add the 4 checksum bytes from stage 7 at the end of extended PIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.")
	bcpy = append(bcpy, hash2[0:4]...)
	fmt.Println(ByteToString(bcpy))
	fmt.Println("================")

	/* Encode base58 string */
	s = B58encode(bcpy)

	/* For number  of leading 0's in bytes, prepend 1 */
	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}
	fmt.Println("8 - Convet the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format")
	fmt.Println(s)
	fmt.Println("================")

	return s
}

func GenerateRstring(l int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//判断是否未编译，使用go run 或者 debug运行
func Isgorun(path string) bool{
	plist:=strings.Split(path,"/")
	if plist[0]=="tmp"{
		return true
	}
	return false
}

//Get current directory
func GetCurrent(path string)string{
	plist:=strings.Split(path,"/")
	current:=strings.Join(plist[:len(plist)-1],"/")
	return current
}