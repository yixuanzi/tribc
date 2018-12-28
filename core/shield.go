// shield.go
package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"strings"
	"tribc/lib"
)


//创建隐私地址
func CreateShieldAddr(acc *Account) ([]byte, []byte ) {
	gkeyA := acc.GkeyA
	gkeyB := acc.GkeyB

	randomkey, _ := ecdsa.GenerateKey(curve, strings.NewReader(lib.GenerateRstring(45)))


	P, _ := curve.ScalarMult(gkeyA.PublicKey.X, gkeyA.PublicKey.Y, randomkey.D.Bytes()) //Mr

	x, y := curve.ScalarBaseMult(P.Bytes()) //(Mr)G

	P, _ = curve.Add(x, y, gkeyB.PublicKey.X, gkeyB.PublicKey.Y) //(Mr)G+N

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

const (
	version            = byte(0x01)
	privKeyBytesLen    = 32
)
