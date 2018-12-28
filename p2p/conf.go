package p2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"tribc/core"
	"tribc/lib"
)

type NodeConf struct {
	Priv string
	Hash string
	Version string

}

type P2PConf struct {
	Btime int
	Lis string
	LRPC string
	RRPC string
	Node []string
}


type Node struct {
	Gkey *core.GKey
	Btime int
	Lis string
	LRPC string
	RRPC string
	Node []string
}
var Curve = elliptic.P256() // 椭圆曲线参数,公共参数

//{"Priv":"92c2d8c8989e90a98e3a1e5c21e7853627cd3c25f27bb33fe057efdf115f5b4e","Hash":"ffa59174556a96061380d9bf2322cd3c","Version":"1.0"}
const PoA  = "8fcd85b226ba9b13cc965e7543eef01e1c228d4722e4457c9f381c8f56bdf94a36f7a9a25705197dc67ebda84a7100190de25bf438223a95b487213ef4b165c2"

var PoAB,_ =hex.DecodeString(PoA)
var PoA_pub = ecdsa.PublicKey{Curve,new(big.Int).SetBytes(PoAB[:32]),new(big.Int).SetBytes(PoAB[32:])}






func Init(path string) bool {
	//gkey,_ := core.MakeNewKey("generate addresss with POA,abcdefghijklmnopqr") //POA secret key

	gkey,_ := core.MakeNewKey(lib.GenerateRstring(45))
	priv :=gkey.GetPrivKey()
	md5key := md5.Sum(priv)
	hash := hex.EncodeToString(md5key[:])
	nc := NodeConf{hex.EncodeToString(priv),hash,"1.0"}
	data,_ := json.Marshal(nc)
	if ioutil.WriteFile(path,data,0644)==nil{
		fmt.Println("[Info Init]","初始化共识节点配置成功！",path)
		return true
	}
	return false
}

var NODE * Node

func Load(consensus string,p2p string) *Node{
	var nc NodeConf
	data,_ := ioutil.ReadFile(consensus)
	json.Unmarshal(data,&nc)
	priv,_:= hex.DecodeString(nc.Priv)
	gkey:=core.Priv2gkey(priv)

	var pc P2PConf
	data,_ = ioutil.ReadFile(p2p)
	json.Unmarshal(data,&pc)
	node:=Node{gkey,pc.Btime,pc.Lis,pc.LRPC,pc.RRPC,pc.Node} //初始化节点参数
	NODE=&node
	return &node
}



