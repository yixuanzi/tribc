package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"encoding/hex"
	"tribc/wasm"
)

type CSA struct {
	Shieldaddr string //`json:shieldaddr`
	ShieldpKey string //`json:shieldpkey`
}

//执行主函数
func main() {
	//android chrome 70
	//android firefox 63
	//ios safair 11
	c := make(chan struct{}, 0)
	//===============================
	//测试代码
	println("Hello, WebAssembly!")
	doc := js.Global().Get("document")

	wasmtest := func(i []js.Value) {
		fmt.Println("this is a wasm test callback func!")
		//fmt.Println(i[0].String())
		//doc.Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
		i[1].Invoke(i[0]) //通过回调函数传递返回计算结果
	}

	doc.Set("wasmtest", js.NewCallback(wasmtest))
	//===============================

	var acc *wasm.Account
	triacc := js.Global().Get("triacc") //使用trias_wasm.wasm的页面必须要创建triacc对象

	//以下创建所有的triacc功能函数变量

	CreateAcc:= func(i []js.Value){ //p1:加密密码 p2:回调函数，并返回账号结果
		handlerError("CreateAcc")
		entry_acc:=wasm.CreateAcc([]byte(i[0].String()))
		i[1].Invoke(js.ValueOf(entry_acc)) //通过回调函数传递返回计算结果
	}

	Load4Data:=func(i []js.Value){ //p1:加密账号数据  p2:解密密码  p3:回调函数，并返回账号结果
		handlerError("Load4Data")
		acc_once:=wasm.Load4Data(i[0].String(),[]byte(i[1].String()))
		if acc_once!=nil{
			acc=acc_once
			i[2].Invoke()
		}
	}

	GetCurrentAcc:=func(i []js.Value){  //p1:回调函数，并返回当前账号地址
		handlerError("GetCurrentAcc")
		if acc!=nil{
			addr:=wasm.GetAddress(acc.GkeyA.GetPubKey(),acc.GkeyB.GetPubKey())
			i[0].Invoke(js.ValueOf(addr))
		}else{
			i[0].Invoke(js.ValueOf("NULL"))
		}
	}

	DestoryCurrAcc:=func(i []js.Value){
		acc=nil
		fmt.Println("Destory Current Acc Succ!")
	}

	Sign:=func(i []js.Value){ //p1:需签名数据  p2:回调函数，并返回签名数据结构，包含公钥和签名数据
		handlerError("Sign")
		if acc!=nil{
			signdata,_:= wasm.Sign(acc.GkeyB.PrivateKey,[]byte(i[0].String()))
			tsign:=wasm.TSign{hex.EncodeToString(acc.GkeyB.GetPubKey()),signdata}
			rs,_:=json.Marshal(tsign)
			i[1].Invoke(js.ValueOf(string(rs)))
		}else{
			i[1].Invoke(js.ValueOf("Fail"))
		}
	}

	Verify:=func(i []js.Value){ //p1:公钥 p2:签名数据 p3:签名验证数据 p4:回调函数，返回验证结果
		handlerError("Verify")
		pub_b,_:= hex.DecodeString(i[0].String())
		pubk:= wasm.Pub2pubKey(pub_b)
		f,_:=wasm.Verify([]byte(i[2].String()),i[1].String(),pubk)
		if f{
			i[3].Invoke(js.ValueOf("OK"))
		}else{
			i[3].Invoke(js.ValueOf("Fail"))
		}

	}


	CreateShieldAddr:= func(i []js.Value) { //p1:用户对外地址 p2:回调函数，返回隐私地址相关数据
		handlerError("CreateShieldAddr")
		if acc!=nil {
			shieldaddr, shieldpKey := wasm.CreateShieldAddr(i[0].String())
			result:=CSA{hex.EncodeToString(shieldaddr),hex.EncodeToString(shieldpKey)}
			rs,_:=json.Marshal(result)
			i[1].Invoke(js.ValueOf(string(rs)))
		}else{
			i[1].Invoke(js.ValueOf("Fail"))
		}
	}

	Verify_shield:=func (i []js.Value) { //p1:shieldaddr p2:shieldpkey p3:回调函数，返回判断状态
		handlerError("Verify_shield")
		saddr, _ := hex.DecodeString(i[0].String())
		spkey, _ := hex.DecodeString(i[1].String())
		if wasm.Verify_shield(acc, saddr, spkey) {
			i[2].Invoke(js.ValueOf("OK"))
		} else {
			i[2].Invoke(js.ValueOf("Fail"))
		}
	}

	Verify_shield2 := func(i []js.Value){ //p1:私钥A p2:账户地址 p3:shieldaddr p4:shieldpkey p5:回调函数，返回判断状态
		handlerError("Verify_shield")
		acc_new := wasm.GetAcc4privA(i[0].String(),i[1].String())
		saddr, _ := hex.DecodeString(i[2].String())
		spkey, _ := hex.DecodeString(i[3].String())

		if wasm.Verify_shield(acc_new, saddr, spkey) {
			i[4].Invoke(js.ValueOf("OK"))
		} else {
			i[4].Invoke(js.ValueOf("Fail"))
		}
	}


	Shield_Sign := func(i []js.Value){ //p1:shieldpkey p2:需要签名数据 p3:回调函数，返回签名数据结构
		handlerError("Shield_Sign")
		if acc != nil {
			spkey,_:=hex.DecodeString(i[0].String())
			priv := wasm.Getprivkey(acc,spkey)
			signdata,_ := wasm.Sign(priv,[]byte(i[1].String()))
			pubKey := append(priv.PublicKey.X.Bytes(), priv.Y.Bytes()...) // []bytes type
			tsign:=wasm.TSign{hex.EncodeToString(pubKey),signdata}
			rs,_:=json.Marshal(tsign)
			i[2].Invoke(js.ValueOf(string(rs)))
		}else {
			i[2].Invoke(js.ValueOf("Fail"))
		}
	}

	GetPrivkeyA := func(i []js.Value){ //p1:回调函数，返回私钥A
		handlerError("GetPrivkeyA")
		if acc != nil{
			privA := acc.GkeyA.GetPrivKey()
			privA_aes,_ := wasm.AesEncrypt(privA,[]byte("tribc"))
			privA_hex := hex.EncodeToString(privA_aes)
			i[0].Invoke(js.ValueOf(privA_hex))
		}else{
			i[0].Invoke(js.ValueOf("Fail"))
		}
	}


	//为对应的triacc函数调用做映射
	triacc.Set("CreateAcc", js.NewCallback(CreateAcc))
	triacc.Set("Load4Data", js.NewCallback(Load4Data))
	triacc.Set("GetCurrentAcc", js.NewCallback(GetCurrentAcc))
	triacc.Set("DestoryCurrAcc", js.NewCallback(DestoryCurrAcc))
	triacc.Set("Sign", js.NewCallback(Sign))
	triacc.Set("Verify", js.NewCallback(Verify))
	triacc.Set("CreateShieldAddr", js.NewCallback(CreateShieldAddr))
	triacc.Set("Verify_shield", js.NewCallback(Verify_shield))
	triacc.Set("Verify_shield2", js.NewCallback(Verify_shield2))
	triacc.Set("Shield_Sign", js.NewCallback(Shield_Sign))
	triacc.Set("GetPrivkeyA", js.NewCallback(GetPrivkeyA))
	<-c
}

//捕获异常处理
func handlerError(name string) {
	if p := recover(); p != nil {
		fmt.Println("[Error] %s call error: %v",name,p)
	}
}