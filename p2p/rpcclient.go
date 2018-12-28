package p2p

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
	"tribc/inc"
)

type CP2PRPC struct {
	conn *rpc.Client
}

//为p2p和钱包提供访问core RPC的快捷调用接口

var CPR *CP2PRPC

func (cpr *CP2PRPC)ConnectServer(addr string){
	for{
		conn, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			fmt.Println("[WARN] Connect RPC Server fail,", err)
		}else{
			cpr.conn=conn
			fmt.Println("[Info] Connect RPC Server succ")
			break
		}
		time.Sleep(time.Duration(3)*time.Second)
	}

}

func (cpr *CP2PRPC)PutEvent(ne *inc.NewEvent,src string)bool{
	var req = inc.NERequest{*ne,src}
	var res inc.RPCResponse
	fmt.Println("[INFO] PutEvent to TRI",req.NE)
	err := cpr.conn.Call("TRIRPC.PutEvent",req,&res)
	if err != nil || res.RS!="OK"{
		log.Fatalln("[ERROR] RPC Call PutEvent to consensus error: ", err)
	}
	return true
}


func (cpr *CP2PRPC)PutTriBlock(tb *inc.TriBlock,src string)bool{
	var req=inc.TBRequest{*tb,src}
	var res inc.RPCResponse
	fmt.Println("[INFO] PutTriBlock to TRI",req.TB)
	err := cpr.conn.Call("TRIRPC.PutTriBlock",req,&res)
	if err != nil || res.RS!="OK"{
		log.Fatalln("[ERROR] RPC Call PutTriBlock to tri error: ", err)
	}
	return true
}


func (cpr *CP2PRPC)GetBlockHeight()(int,string){
	req := inc.CMDRequest{"GetBlockHeight"}
	var res inc.HeightResponse
	err := cpr.conn.Call("TRIRPC.GetBlockHeight",req,&res)
	if err != nil {
		fmt.Println("[ERROR] RPC Call GetBlockHeight error: ", err)
		return -2,""
	}
	return res.Height,res.Hash
}

func (cpr *CP2PRPC)GetUTXO(utxo string)* inc.UTXO{
	var req = inc.CMDRequest{utxo}
	var res inc.UTXOResponse
	err := cpr.conn.Call("TRIRPC.GetUTXO",req,&res)
	if err != nil {
		fmt.Println("[ERROR] RPC Call GetUTXO error: ", err)
		return nil
	}
	return &res.UOBJ
}

func Initcp2prpc(addr string)bool{
	CPR=new(CP2PRPC)
	go CPR.ConnectServer(addr)
	return true
}
