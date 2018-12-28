package swdb

import (
	"encoding/json"
	"fmt"
	"tribc/inc"
)



func printErr(err error,name string){
	if err != nil {
		fmt.Printf("[Error %s]:%s\n",name,err)
	}
}


var SW *StatusWorld

type StatusWorld struct {
	t triasdb
}


func InitSW(path string)bool{
	var sw StatusWorld
	sw.Init_sw(path)
	SW=&sw
	return true
}

func (sw *StatusWorld) Init_sw(path string) bool {
	f,err:=sw.t.initswdb(path)
	if err!=nil{
		fmt.Println("[Error] Init_sw fail",err)
	}else{
		fmt.Println("[Info] init_sw succ!")
	}
	return f
}
func (sw *StatusWorld) UTXO_add(hash string, utxo *inc.UTXO) bool {
	data,_:= json.Marshal(utxo)
	return sw.t.set([]byte(hash), data)

}

func (sw *StatusWorld) UTXO_del(hash string) bool {
	return sw.t.del([]byte(hash))
}

func (sw *StatusWorld) UTXO_get(hash string) *inc.UTXO {
	data:=sw.t.get([]byte(hash))
	if len(data)==0{
		return nil
	}
	var utxo inc.UTXO
	json.Unmarshal(data,&utxo)
	return &utxo
}

func (sw *StatusWorld) Block_add(hash string,sb *inc.Sblock) bool{
	data,_:= json.Marshal(sb)
	return sw.t.set([]byte(hash),data)
}

func (sw *StatusWorld) Event_add(hash string,ne *inc.Event) bool{
	data,_:= json.Marshal(ne)
	return sw.t.set([]byte(hash),data)
}

func (sw *StatusWorld) Block_get(hash string) *inc.Sblock{
	data:= sw.t.get([]byte(hash))
	if len(data)==0{
		return nil
	}
	var sb inc.Sblock
	json.Unmarshal(data,&sb)
	return &sb
}

func (sw *StatusWorld) Event_get(hash string) *inc.Event{
	data:= sw.t.get([]byte(hash))
	if len(data)==0{
		return nil
	}
	var ee inc.Event
	json.Unmarshal(data,&ee)
	return &ee
}


func (sw *StatusWorld)Account_get(hash string) * inc.AccountStatus{
	data:=sw.t.get([]byte(hash))
	if len(data)==0{
		return nil
	}
	var as inc.AccountStatus
	json.Unmarshal(data,&as)
	return &as
}

func (sw *StatusWorld)Account_set(hash string,as *inc.AccountStatus)bool{
	data,_:= json.Marshal(as)
	return sw.t.set([]byte(hash),data)
}





func (sw *StatusWorld) Close_sw()bool{
	sw.t.close()
	return true
}