package core

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
	"tribc/inc"
	"tribc/swdb"
)

// rpc操作对象结构体
type TRIRPC struct {
}


// 由下层共识模块传递而来的事件，对其进行业务合法性检查之后再传递到共识模块
func (this *TRIRPC) PutEvent(req inc.NERequest, res *inc.RPCResponse) error {
	if CheckNewEvent(&req.NE){
		fmt.Println("Have PutEvent:",req.NE)
		time.Sleep(time.Duration(3)*time.Second)
		res.RS="OK"
		var rs inc.RPCResponse
		go CR.PutEvent(req,&rs) //判断事件合法，再次传递到共识模块中
	}
	return nil
}

func (this *TRIRPC) PutTriBlock(req inc.TBRequest,res *inc.RPCResponse) error{
	fmt.Println("Have PutTriBlock:",req.TB)
	go ProcTriBlock(&req.TB,swdb.SW)
	res.RS="OK"
	return nil
}

func (this *TRIRPC) GetUTXO(req inc.CMDRequest,res *inc.UTXOResponse) error {
	utxo:= swdb.SW.UTXO_get(req.Cmd)
	res.UOBJ=*utxo
	return nil


}

func (this *TRIRPC) GetBlockHeight(req inc.CMDRequest,res *inc.HeightResponse) error{
	res.Height=Height
	res.Hash=CurentHash
	return nil
}


var CC *inc.CoreConf

func CreateRPC(cc * inc.CoreConf) bool{
	rpc.Register(new(TRIRPC)) // 注册rpc服务
	rpc.HandleHTTP()         // 采用http协议作为rpc载体

	lis, err := net.Listen("tcp", cc.LRPC)
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}
	fmt.Println("CreateRPC succ",cc.LRPC)
	go http.Serve(lis, nil)
	return true
}
