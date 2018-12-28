package p2p

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"tribc/core"
	"tribc/inc"
)


//combase 权威地址：4TSbzS4t3NSVw4aEid93rApbUEMX
func (pn *P2PNet)Generateblock(){
	if pn.node.Btime==0{
		return
	}
	index:= 0
	for{
		time.Sleep(time.Duration(30)*time.Second)
		if CPR.conn==nil{ //若未连接到Trias应用，不产生块
			continue
		}
		var evlist []inc.NewEvent
		ccoin:= pn.CreateCoin("4TSbzS4t3NSVw4aEid93rApbUEMX")
		evlist=append(evlist,*ccoin)

		if len(pn.eventcache)>10{
			index=10
		}else{
			index=len(pn.eventcache)
		}

		for i:=0;i<index;i++{
			evlist=append(evlist,pn.eventcache[i])
		}
		pn.eventcache=pn.eventcache[index:]
		p_height,p_hash:= CPR.GetBlockHeight() //基于RPC获得当前区块链最顶端的快高和hash
		tb:=inc.TriBlock{"",p_hash,p_height+1,"adminMiner",inc.TSign{"",""},"None","None",time.Now().Unix(),evlist}
		data,_:=json.Marshal(tb)
		hash,tsign:=core.ComputerHashSign(data,NODE.Gkey)
		tb.Hash=hash
		tb.Sign=*tsign


		fmt.Println("Generate Triblock ",tb)
		data,_=json.Marshal(tb)
		pn.evented[tb.Hash]=tb.Timestamp //上传块到trias，以执行块
		pn.Broadcast(data,"127.0.0.1",2) //广播块
		CPR.PutTriBlock(&tb,"127.0.0.1")
	}
}

//铸币交易
func (pn *P2PNet)CreateCoin(addr string) *inc.NewEvent{
	ne := inc.NewEvent{"",0,0,"TRI",inc.TSign{"",""},time.Now().Unix(),"Trias Create Coin", addr+"|10"}
	data,_:=json.Marshal(ne)
	hash,_ := core.ComputerHashSign(data,nil)
	ne.Exhash=hash
	return &ne
}

func VerifyBlock(tb *inc.TriBlock,pub *ecdsa.PublicKey)bool{
	_,p_hash:= CPR.GetBlockHeight()
	if tb.ParentHash==p_hash {
		hash,_ := hex.DecodeString(tb.Hash)
		if ok, _ := core.Verify(hash, tb.Sign.Sigadata, pub);ok{
			return true
		}
	}
	return false
}

