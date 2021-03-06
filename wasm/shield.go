// shield.go
package wasm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"strings"
)


//创建隐私地址
func CreateShieldAddr(addr string) ([]byte, []byte ) {
	//gkeyA := acc.GkeyA
	//gkeyB := acc.GkeyB
	pubA,pubB:=GetPubk4Addr(addr)

	A_X:=pubA.X
	A_Y:=pubA.Y
	B_X:=pubB.X
	B_Y:=pubB.Y
	randomkey, _ := ecdsa.GenerateKey(curve, strings.NewReader(GenerateRstring(45)))


	P, _ := curve.ScalarMult(A_X, A_Y, randomkey.D.Bytes()) //Mr

	x, y := curve.ScalarBaseMult(P.Bytes()) //(Mr)G

	P, _ = curve.Add(x, y, B_X, B_Y) //(Mr)G+N

	pubkey := append(randomkey.PublicKey.X.Bytes(),randomkey.PublicKey.Y.Bytes()...)
	return P.Bytes(),pubkey
}


//验证隐私地址
func Verify_shield(acc *Account, shieldaddr []byte, shieldpKey []byte) bool {
	gkeyA := acc.GkeyA
	gkeyB := acc.GkeyB

	P := new(big.Int).SetBytes(shieldaddr)

	R_x := new(big.Int).SetBytes(shieldpKey[:32])
	R_y := new(big.Int).SetBytes(shieldpKey[32:])

	p, _ := curve.ScalarMult(R_x, R_y, gkeyA.PrivateKey.D.Bytes()) //mR

	x, y := curve.ScalarBaseMult(p.Bytes()) //(mR)G

	p, _ = curve.Add(x, y, gkeyB.PublicKey.X, gkeyB.PublicKey.Y) //(mR)G+N

	if P.Cmp(p) == 0 {
		return true
	}
	return false
}

//根据隐私地址，随机数，账号私钥，计算出隐私地址私钥
func Getprivkey(acc *Account, shieldpKey []byte) *ecdsa.PrivateKey {
	gkeyA := acc.GkeyA
	gkeyB := acc.GkeyB

	R_x := new(big.Int).SetBytes(shieldpKey[:32])
	R_y := new(big.Int).SetBytes(shieldpKey[32:])

	x, _ := curve.ScalarMult(R_x, R_y, gkeyA.PrivateKey.D.Bytes()) //mR
	x = x.Add(x, gkeyB.PrivateKey.D)                               //(mR+n)

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = x
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(x.Bytes()) //(mR+n)G
	return priv
}


var curve = elliptic.P256() // 椭圆曲线参数,公共参数
var N = curve.Params().N


