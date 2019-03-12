package main

import (
	"fmt"
	"strings"
	"tribc/wasm"
	"crypto/ecdsa"
	"crypto/elliptic"
)

var curve = elliptic.P256() // 椭圆曲线参数,公共参数
var N = curve.Params().N

func test(rs string,acc *wasm.Account,addr string){
	shieldaddr, shieldpKey := wasm.CreateShieldAddr(addr,rs)

	if wasm.Verify_shield(acc,shieldaddr,shieldpKey) {
		fmt.Println("OK")
	}else{
		fmt.Println("Fail")
	}

	pubk:=wasm.GetPubk4Addr(addr)
	A_X:=pubk.X
	A_Y:=pubk.Y

	randomkey, _ := ecdsa.GenerateKey(curve, strings.NewReader(rs))

	P, _ := curve.ScalarMult(A_X, A_Y, randomkey.D.Bytes()) //Mr
	//P.Mod(P,N)

	p, _ := curve.ScalarMult(randomkey.X, randomkey.Y, acc.GkeyA.PrivateKey.D.Bytes()) //mR
	//p.Mod(p,N)

	if P.Cmp(p)==0{
		fmt.Println("OKKK")
	}else{
		fmt.Println("FAILL")
	}
}

func main()  {

	entry_acc:=wasm.CreateAcc([]byte("1234qwer"),"wm82hl7o1ffklx8wi8yxbe35je55gr46muh4p0j7mk0dm")
	acc:=wasm.Load4Data(entry_acc,[]byte("1234qwer"))
	addr:=wasm.GetAddress(acc.GkeyA.GetPubKey())

	test("w3e74ld2mply80smof1t8b7p19oiz10o2os9bzexg1ovt",acc,addr)
	test("a3e74ld2mply80smof1t8b7p19oiz10o2os9bzexg1ovt",acc,addr)
	test("gvzeyligvrvmrq7hik7rmzm7cnbl15gexkl1k7zzyavuz",acc,addr)
	test("vwoqcqxt033nxqstw2ypm0jj5bcvfx60442gksdt5h1j2",acc,addr)


}