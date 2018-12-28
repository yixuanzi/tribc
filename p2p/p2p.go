package p2p

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"tribc/inc"
)

type P2PNet struct {
	node *Node
	p2pnode []Dstp2p
	evented map[string]int64 //存储已处理的事件,当前此对象处于多线程非安全状态，后期改进
	eventcache []inc.NewEvent //缓存的有效新事件
}

var PN *P2PNet

func (pn *P2PNet)CreateServer()bool{
	go pn.CreateServerRouter()
	return true
}

func (pn *P2PNet)ConnServer()bool{
	for i:=0;i<len(pn.p2pnode);i++{
		go pn.p2pnode[i].ConnServer()
	}
	return true
}

//从服务端进行tcp整流和消息完整性检查后传递过来的消息结构体，后续需要对广播过来的消息进行一系列业务处理，数据主要包括事件和块数据
func (pn *P2PNet)HandleMsg(mess *Msg){
	//此处接收到的消息来自其他节点广播来数据，根据不同数据进行不同处理
	//只处理最近10分钟的事件
	//事件数据：根据事件hash检验当前节点是否处理过，若没有则传递到业务rpc进行业务合规性检查
	//fmt.Println(mess)
	//广播过来的数据结构为：type-data    mess.Mess[0]=type mess.Mess[1:]=data
	now := time.Now().Unix()
	if mess.Mess[0]==1{ //事件数据
		fmt.Println("[INFO] Get Event...")
		var ne inc.NewEvent
		json.Unmarshal(mess.Mess[1:],&ne)
		if now - ne.Timestamp > 600{
			return //不处理10分钟之前的数据
		}
		if _,ok := pn.evented[ne.Exhash];!ok{
			//rcp 传递数据
			CPR.PutEvent(&ne,mess.Addr)
			pn.evented[ne.Exhash]=ne.Timestamp
		}
	}else if mess.Mess[0]==2{ //块数据
		fmt.Println("[INFO] Get Block...")
		var tb inc.TriBlock
		json.Unmarshal(mess.Mess[1:],&tb)
		if now - tb.Timestamp > 600{
			return //不处理10分钟之前的数据
		}
		if _,ok := pn.evented[tb.Hash];!ok{
			//PoA共识（校验权威节点签名,判断前后块是否连续）
			if !VerifyBlock(&tb,&PoA_pub){
				return
			}
			//rcp 传递数据
			CPR.PutTriBlock(&tb,mess.Addr)
			pn.evented[tb.Hash]=tb.Timestamp
		}
	}else {
		fmt.Println("[WARN] 无效的消息数据",mess.Addr,string(mess.Mess))
	}

}


//广播数据
func (pn *P2PNet)Broadcast(mess []byte,Src string,msgtype uint8) bool{
	for i:=0;i<len(pn.p2pnode);i++{
		//if pn.p2pnode[i].IP!=Src{ //广播时对数据来源服务不广播
			pn.p2pnode[i].Sendmsg(mess,msgtype)
		//}
	}
	return true
}

func CreateNet(node *Node)* P2PNet{
	var pn P2PNet

	pn.node=node
	pn.evented=make(map[string]int64)
	PN=&pn

	go pn.clearevented() //周期性的清理已处理事件map

	for _,nodeaddr := range (node.Node){
		if len(nodeaddr)>5{
			pnode := Dstp2p{nodeaddr,nil,nil,strings.Split(nodeaddr,":")[0]}
			pn.p2pnode=append(pn.p2pnode,pnode)
			//append(p2p.p2pnode,p2n)
			//fmt.Println(p2n)
		}
	}
	return &pn
}


func (pn *P2PNet)clearevented(){
	now := time.Now().Unix()
	for{
		time.Sleep(time.Duration(10)*time.Minute) //10分钟的清理周期
		now = time.Now().Unix()

		for hash:= range(pn.evented){
			if (now - pn.evented[hash] > 60*10){ //清理缓存了10分钟之外的数据
				delete(pn.evented,hash)
			}
		}
	}
}

