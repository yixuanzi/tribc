package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"tribc/inc"
)

// rpc操作对象结构体
type P2PRPC struct {
}


// 传递事件到共识层，由core（上层应用）调用处理
func (this *P2PRPC) PutEvent(req inc.NERequest, res *inc.RPCResponse) error {
	//此处传递来的事件属于当前节点处理的事件，主要完成 事件缓存（后期块生成）和 事件广播的功能
	fmt.Println("Get Checked NewEvent from Tri ",req.NE)
	PN.eventcache=append(PN.eventcache,req.NE) //添加缓存
	if _,ok := PN.evented[req.NE.Exhash];!ok{    //若是本节点产生的事件，添加处理记录
		PN.evented[req.NE.Exhash]=req.NE.Timestamp
	}
	data,_:=json.Marshal(req.NE)
	res.RS="OK"
	go PN.Broadcast(data,req.Src,1) //广播事件
	return nil
}

func CreateRPC(node *Node) bool{
	rpc.Register(new(P2PRPC)) // 注册rpc服务
	rpc.HandleHTTP()         // 采用http协议作为rpc载体

	lis, err := net.Listen("tcp", node.LRPC)
	if err != nil {
		fmt.Println("Create RPC Server fail", err)
		os.Exit(2)
	}
	fmt.Println("CreateRPC Server succ",node.LRPC)
	go http.Serve(lis, nil)
	return true
}