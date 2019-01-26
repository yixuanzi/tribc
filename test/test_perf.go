package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
	"tribc/core"
	"tribc/lib"
)

func main() {
	//生成账号测试
	fmt.Println("Generate addresss for account")
	//创建PoA权威账号
	//gkeyA, err := core.MakeNewKey("123456789012345678901234567890123456789012345")
	//gkeyB, err := core.MakeNewKey("098765432109876543210987654321098765432154321")

	randomstr:=lib.GenerateRstring(45)
	acc:=*core.CreateAccount(randomstr)
	addr := core.GetAddress(acc.GkeyA.GetPubKey())

	fmt.Println("账号转化地址为:",addr)
	fmt.Println("==========================")
	//隐私地址功能测试
	shieldaddr, shieldpKey := core.CreateShieldAddr(addr)
	/*
	s_1:=lib.ByteToString(shieldaddr)
	s_1 = hex.EncodeToString(shieldaddr)
	s_2,_:=hex.DecodeString(s_1)
	fmt.Println(s_2)
	*/
	fmt.Println("A&B shieldaddr is :", lib.ByteToString(shieldaddr))
	fmt.Println("A&B shieldpKey is :", lib.ByteToString(shieldpKey))

	statzq:=int64(5)
	//analysis Verify_shield
	stat:=int64(0)
	now:= time.Now().Unix()
	for{
		core.Verify_shield(&acc, shieldaddr, shieldpKey)
		stat++
		if time.Now().Unix()-now > statzq{
			break
		}
	}
	fmt.Println("The core.Verify_shield performance: ",stat/statzq)



	//Virify(gkeyA,gkeyB,b_s_1,b_s_2)
	if core.Verify_shield(&acc, shieldaddr, shieldpKey) {
		//println("The virify is succfully!")
	}

	priv := core.Getprivkey(&acc, shieldpKey)
	text := []byte("hahahaha~!")
	r, s, _ := ecdsa.Sign(rand.Reader, priv, text)
	if ecdsa.Verify(&priv.PublicKey, text, r, s) {
		//println("The sign virify is succfully!")
	}


	stat=int64(0)
	now = time.Now().Unix()
	for{
		_,_ = core.Sign(priv,text)
		stat++
		if time.Now().Unix()-now > statzq{
			break
		}
	}
	fmt.Println("The core.Sign performance: ",stat/statzq)


	stat = int64(0)
	now = time.Now().Unix()
	stext,_ := core.Sign(priv,text)
	for{
		_,_=core.Verify(text,stext,&priv.PublicKey)
		stat++
		if time.Now().Unix()-now > statzq{
			break
		}
	}
	fmt.Println("The core.Verify performance: ",stat/statzq)


	//账号加解密测试
	var aeskey = []byte("1234qwer1234qwer") //长度必须为16,24,32
	//pass := []byte("vdncloud123456")
	pass:=[]byte(lib.GenerateRstring(1024*1024))

	stat = int64(0)
	now = time.Now().Unix()
	for{
		_, _ = lib.AesEncrypt(pass, aeskey)
		stat++
		if time.Now().Unix()-now > statzq{
			break
		}
	}
	fmt.Println("The lib.AesEncrypt performance: ",stat/statzq)


	xpass, _ := lib.AesEncrypt(pass, aeskey)
	pass64 := base64.StdEncoding.EncodeToString(xpass)
	//fmt.Printf("加密后:%v\n",pass64)
	bytesPass, _ := base64.StdEncoding.DecodeString(pass64)

	stat = int64(0)
	now = time.Now().Unix()
	for{
		_, _ = lib.AesDecrypt(bytesPass, []byte("1234qwer1234qwer"))
		stat++
		if time.Now().Unix()-now > statzq{
			break
		}
	}
	fmt.Println("The lib.AesDecrypt performance: ",stat/statzq)



}
