package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"tribc/core"
	"tribc/p2p"
)

var Curve = elliptic.P256() // 椭圆曲线参数,公共参数

const PoA  = "8fcd85b226ba9b13cc965e7543eef01e1c228d4722e4457c9f381c8f56bdf94a36f7a9a25705197dc67ebda84a7100190de25bf438223a95b487213ef4b165c2"



var current="/home/lab8/go/src/trias/"

func main(){
	//gkey,_ := core.MakeNewKey("generate addresss with POA,abcdefghijklmnopqr")
	node:= p2p.Load(current+"conf/consensus.conf",current+"conf/p2p.conf")



	PoAA,_:=hex.DecodeString(PoA)

	var PoA_pub = ecdsa.PublicKey{Curve,new(big.Int).SetBytes(PoAA[:32]),new(big.Int).SetBytes(PoAA[32:])}

	h:=sha256.New()
	h.Write([]byte("abcdefg"))
	hash:= h.Sum(nil)
	fmt.Println(hex.EncodeToString(hash))

	signdata,_:= core.Sign(node.Gkey.PrivateKey,hash)

	flag,_ := core.Verify(hash,signdata,&PoA_pub)
	fmt.Println("OK",flag)

	signdata,_= core.Sign(node.Gkey.PrivateKey,hash)
	flag,_ = core.Verify(hash,signdata,&PoA_pub)
	fmt.Println("OK",flag)

	ts,tsign := core.ComputerHashSign([]byte("abcdefg"),node.Gkey)
	hash,_= hex.DecodeString(ts)
	flag,_ = core.Verify(hash,tsign.Sigadata, &PoA_pub)

	if hex.EncodeToString(hash)==ts{
		fmt.Println("OOOOOOOOOOKKKKKKKK1")
	}

	if tsign.Pubkey==PoA{
		fmt.Println("OOOOOOOOOOKKKKKKKK2")
	}


	fmt.Println("OK",flag,ts)
}
