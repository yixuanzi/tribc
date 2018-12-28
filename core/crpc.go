package core

import (
	"fmt"
	"net/rpc"
	"time"
	"tribc/inc"
)

type CRPC struct {
	conn *rpc.Client
}

var CR *CRPC

func (cr *CRPC)ConnectServer(addr string){
	for{
		conn, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			fmt.Println("[ERROR] Connect RPC Server fail", err)
		}else{
			cr.conn=conn
			fmt.Println("[Info] Connect RPC Server succ", err)
			break
		}
		time.Sleep(time.Duration(3)*time.Second)
	}
}

func (cr *CRPC)PutEvent(req inc.NERequest, res *inc.RPCResponse)bool{
	fmt.Println("[INFO] PutEvent ",req.NE)
	err := cr.conn.Call("P2PRPC.PutEvent",req,res)
	if (err != nil) || (res.RS!="OK") {
		fmt.Println("[ERROR] RPC Call PutEvent to consensus error: ", err)
	}
	return true
}

func InitCRPC(cc *inc.CoreConf)bool{
	CR=new(CRPC)
	go CR.ConnectServer(cc.RRPC)
	return true
}

