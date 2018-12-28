package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"tribc/inc"
	"tribc/swdb"
)

var Height=-1
var CurentHash="0000000000000000000000000000000000000000000000000000000000000000"

func CheckNewEvent(ne * inc.NewEvent)bool{
	if ne.EventType==0{
		switch ne.EventID {
		case 0: return Fund0_event_check(ne)
		case 1: return Fund1_event_check(ne)
		case 2: return Fund2_event_check(ne)
		case 3: return Fund3_event_check(ne)
		case 4: return Fund4_event_check(ne)
		default:
			return false
		}
	}
	//目前只处理资金事件
	return false
}


func ProcNewEvent(ne *inc.NewEvent,height int,sw *swdb.StatusWorld)bool{
	if ne.EventType==0{
		switch ne.EventID {
		case 0: return Fund0_event_proc(ne,height)
		case 1: return Fund1_event_proc(ne,height)
		case 2: return Fund2_event_proc(ne,height)
		case 3: return Fund3_event_proc(ne,height)
		case 4: return Fund4_event_proc(ne,height)
		default:
			return false
		}
	}
	//目前只处理资金事件
	return false
}


func ProcTriBlock(tb *inc.TriBlock,sw *swdb.StatusWorld)bool{
	if tb.ParentHash!=CurentHash{ //目前仅进行前后哈希链检查（
		return false
	}
	var sb inc.Sblock
	sb.Hash=tb.Hash
	sb.ParentHash=tb.ParentHash
	sb.Height=tb.Height
	sb.Miner=tb.Miner
	sb.Sign=tb.Sign
	sb.StatusRoot=tb.StatusRoot
	sb.EventRoot=tb.EventRoot
	sb.Timestamp=tb.Timestamp
	for i:=0;i<len(tb.Eventlist);i++{
		sb.Eventlisthash=append(sb.Eventlisthash,tb.Eventlist[i].Exhash)
		ProcNewEvent(&tb.Eventlist[i],tb.Height,sw)
	}
	if sw.Block_add(tb.Hash,&sb){
		CurentHash=tb.Hash //设置当前块
		Height=tb.Height
	}else {
		return false
	}
	return true
}


func default_proc()bool{
	return true
}

//计算数据的哈希和签名,当gkey为nil时，不对hash进行签名
func ComputerHashSign(d []byte,gkey *GKey)(string,*inc.TSign){
	h:=sha256.New()
	h.Write(d)
	hash:= h.Sum(nil)
	hash_s := hex.EncodeToString(hash)
	if gkey!=nil{
		signdata,_:= Sign(gkey.PrivateKey,hash)
		tsign:=inc.TSign{hex.EncodeToString(gkey.GetPubKey()),signdata}
		return hash_s,&tsign
	}
	return hash_s,nil
}
//计算数据的哈希
func ComputerHash(d []byte) []byte{
	h:=sha256.New()
	h.Write(d)
	hash:= h.Sum(nil)
	return hash
}


func TestEvent(){
	for{
		if CR.conn!=nil{
			ne:=inc.NewEvent{"1234567",0,1,"",inc.TSign{"",""},time.Now().Unix(),"this is a input data","out "}
			fmt.Println("PutTestEvent",ne)
			var req=inc.NERequest{ne,"127.0.0.1"}
			var res inc.RPCResponse
			CR.PutEvent(req,&res)
		}else {
			time.Sleep(time.Duration(3)*time.Second)
			fmt.Println("CRPC is not connect,wait again!")
		}
		time.Sleep(time.Duration(3)*time.Second)
	}
}