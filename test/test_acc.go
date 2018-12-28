package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"tribc/core"
	"tribc/lib"
)

func main() {
	//生成账号测试
	fmt.Println("Generate addresss for account")
	//创建PoA权威账号
	//gkeyA, err := core.MakeNewKey("123456789012345678901234567890123456789012345")
	//gkeyB, err := core.MakeNewKey("098765432109876543210987654321098765432154321")

	gkeyA, err := core.MakeNewKey(lib.GenerateRstring(45))
	gkeyB, err := core.MakeNewKey(lib.GenerateRstring(45))
	if err != nil {
		fmt.Println(err)
		return
	}
	privKeyA := gkeyA.GetPrivKey()
	privKeyB := gkeyB.GetPrivKey()
	fmt.Println("B privateKey is :", lib.ByteToString(privKeyA))
	//fmt.Println("A privateKey is :",  hex.EncodeToString(privKeyA)) lib.ByteToString hex.EncodeToString 为等效功能函数
	fmt.Println("B privateKey is :", lib.ByteToString(privKeyB))
	pubKeyA := gkeyA.GetPubKey()
	pubKeyB := gkeyB.GetPubKey()
	fmt.Println("A publickKey is :", lib.ByteToString(pubKeyA))
	fmt.Println("B publickKey is :", lib.ByteToString(pubKeyB))

	acc := core.Account{gkeyA,gkeyB}
	fmt.Println("账号转化地址为:",core.GetAddress(acc.GkeyA.GetPubKey()))
	fmt.Println("==========================")
	//隐私地址功能测试
	shieldaddr, shieldpKey := core.CreateShieldAddr(&acc)
	/*
	s_1:=lib.ByteToString(shieldaddr)
	s_1 = hex.EncodeToString(shieldaddr)
	s_2,_:=hex.DecodeString(s_1)
	fmt.Println(s_2)
	*/
	fmt.Println("A&B shieldaddr is :", lib.ByteToString(shieldaddr))
	fmt.Println("A&B shieldpKey is :", lib.ByteToString(shieldpKey))

	//Virify(gkeyA,gkeyB,b_s_1,b_s_2)
	if core.Verify_shield(&acc, shieldaddr, shieldpKey) {
		println("The virify is succfully!")
	}

	priv := core.Getprivkey(&acc, shieldpKey)
	text := []byte("hahahaha~!")
	r, s, err := ecdsa.Sign(rand.Reader, priv, text)

	if ecdsa.Verify(&priv.PublicKey, text, r, s) {
		println("The sign virify is succfully!")
	}

	stext,_ := core.Sign(priv,text)
	fmt.Println("The sign data:",stext)
	f,err:=core.Verify(text,stext,&priv.PublicKey)
	if f{
		fmt.Println("The sign virfity is scuufully with lib func")
	}

	fmt.Println("==========================")
	//账号加解密测试
	var aeskey = []byte("1234qwer1234qwer") //长度必须为16,24,32
	pass := []byte("vdncloud123456")
	xpass, err := lib.AesEncrypt(pass, aeskey)

	if err != nil {
		fmt.Println(err)
		return
	}

	pass64 := base64.StdEncoding.EncodeToString(xpass)
	fmt.Printf("加密后:%v\n",pass64)

	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	if err != nil {
		fmt.Println(err)
		return
	}

	tpass, err := lib.AesDecrypt(bytesPass, []byte("1234qwer1234qwer"))
	if tpass==nil || err != nil {
		fmt.Println("解密失败",err)
	}
	fmt.Printf("解密后:%s\n", tpass)

	fmt.Println("==========================")
	//账号加密导出到文件
	fmt.Println("正在导出账号文件..........")
	af :=core.Account{gkeyA,gkeyB}
	core.Save2file(&af,"/tmp/trias_acc.json",[]byte("1234qwer"))
	fmt.Println("正在导入账号文件..........")
	facc:= core.Load4file("/tmp/trias_acc.json",[]byte("1234qwer"))
	fmt.Println(facc)



}