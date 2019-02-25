package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tribc/inc"
	"tribc/swdb"
)

/*
(各个参数已经结构化，文本结构无效）
资产 0：铸币事件，只有输出，输出结构为“addr|bal”
	 1： utxo->utxo
     2: acc->acc
	 3: acc->utxo
	 4: utxo->acc

 */

func splitiostring(ios string) []string{
	return  strings.Split(ios,"\n")
}

func splitpstring(ps string) []string{
	return  strings.Split(ps,"|")
}

func Fund0_event_check(ne *inc.NewEvent)bool{
	defer handlerError(ne)
	strs:= splitiostring(ne.OutData)
	ps:= splitpstring(strs[0])
	if ps[1]=="10"{
		return true
	}
	return false
}


func Fund1_event_check(ne *inc.NewEvent)bool{ //整体数据hash校验未进行

	var input inc.UTXO_INPUTS
	json.Unmarshal([]byte(ne.InputData),&input)
	var in_funds float64

	//校验校验输入，无论那个检验出错，都反回假
	for i:=0;i<len(input.UTXOI);i++{
		utxo := swdb.SW.UTXO_get(input.UTXOI[i].Utxo)
		if utxo==nil{
			fmt.Println("[WARN] A invaild utxo")
			return false //不存在的utxo
		}
		pub,_:=hex.DecodeString(input.UTXOI[i].Sign.Pubkey)
		if GetAddress(pub)!=utxo.Owner{
			fmt.Println("[WARN] illegal owner")
			return false //非法拥有者
		}

		hash := ComputerHash([]byte(input.UTXOI[i].Utxo))
		pubkey := Pub2pubKey(pub)

		if f,_ := Verify(hash,input.UTXOI[i].Sign.Sigadata,pubkey);f{
			in_funds+= utxo.Balance
		}else{
			fmt.Println("[WARN] A invaild signdata with event")
			return false //非法签名
		}

	}

	//循环输出，获取输出金额
	var output inc.UTXO_OUTPUTS
	json.Unmarshal([]byte(ne.OutData),&output)
	var out_funds float64
	for i:=0;i>len(output.UTXOO);i++{
		out_funds+= output.UTXOO[i].Balance
	}

	if in_funds < out_funds {
		fmt.Println("[WARN] The event balance is not balance")
		return false
	}
	return true
}

func Fund2_event_check(ne *inc.NewEvent)bool{ //nonce 防重放暂不实现(通过input字段传递nonce；整体数据hash校验未进行
	pub,_:= hex.DecodeString(ne.Sign.Pubkey)
	pubkey := Pub2pubKey(pub)
	hash,_ := hex.DecodeString(ne.Exhash)
	if f,_:= Verify(hash,ne.Sign.Sigadata,pubkey);!f{
		fmt.Println("[WARN] A invaild signdata with event")
		return false
	}

	addr:= GetAddress(pub)
	acc_status:=swdb.SW.Account_get(addr)
	if acc_status==nil{
		fmt.Println("[WARN] A invaild Account")
		return false
	}

	var output inc.ACC_PARAS
	json.Unmarshal([]byte(ne.OutData),&output)
	var out_funds float64
	for i:=0;i<len(output.ACCP);i++{
		out_funds+=output.ACCP[i].Bal
	}

	if acc_status.Bal < out_funds{
		fmt.Println("[WARN] The event balance is not balance")
		return false
	}
	/*
	if acc_status.Nonce+1 == input.nonce{
		ok
	}
	*/
	return true
}

func Fund3_event_check(ne *inc.NewEvent)bool{
	pub,_:= hex.DecodeString(ne.Sign.Pubkey)
	pubkey := Pub2pubKey(pub)
	hash,_ := hex.DecodeString(ne.Exhash)
	if f,_:= Verify(hash,ne.Sign.Sigadata,pubkey);!f{
		fmt.Println("[WARN] A invaild signdata with event")
		return false
	}

	addr:= GetAddress(pub)
	acc_status:=swdb.SW.Account_get(addr)
	if acc_status==nil{
		fmt.Println("[WARN] A invaild Account")
		return false
	}

	//循环输出，获取输出金额
	var output inc.UTXO_OUTPUTS
	json.Unmarshal([]byte(ne.OutData),&output)
	var out_funds float64
	for i:=0;i>len(output.UTXOO);i++{
		out_funds+= output.UTXOO[i].Balance
	}

	if acc_status.Bal < out_funds{
		fmt.Println("[WARN] The event balance is not balance")
		return false
	}

	return true
}

func Fund4_event_check(ne *inc.NewEvent)bool{

	var input inc.UTXO_INPUTS
	json.Unmarshal([]byte(ne.InputData),&input)
	var in_funds float64

	//校验校验输入，无论那个检验出错，都反回假
	for i:=0;i<len(input.UTXOI);i++{
		utxo := swdb.SW.UTXO_get(input.UTXOI[i].Utxo)
		if utxo==nil{
			fmt.Println("[WARN] A invaild utxo")
			return false //不存在的utxo
		}
		pub,_:=hex.DecodeString(input.UTXOI[i].Sign.Pubkey)
		if GetAddress(pub)!=utxo.Owner{
			fmt.Println("[WARN] illegal owner")
			return false //非法拥有者
		}

		hash := ComputerHash([]byte(input.UTXOI[i].Utxo))
		pubkey := Pub2pubKey(pub)

		if f,_ := Verify(hash,input.UTXOI[i].Sign.Sigadata,pubkey);f{
			in_funds+= utxo.Balance
		}else{
			fmt.Println("[WARN] A invaild signdata with event")
			return false //非法签名
		}

	}

	var output inc.ACC_PARAS
	json.Unmarshal([]byte(ne.OutData),&output)
	var out_funds float64
	for i:=0;i<len(output.ACCP);i++{
		out_funds+=output.ACCP[i].Bal
	}

	if in_funds < out_funds{
		fmt.Println("[WARN] The event balance is not balance")
		return false
	}

	return true
}

//==============

func Fund0_event_proc(ne *inc.NewEvent,height int)bool{
	ee := inc.Event{ne.Exhash,ne.EventType,ne.EventID,ne.From,ne.Sign,ne.Timestamp,ne.InputData,ne.OutData,1,height}
	if Fund0_event_check(ne){
		strs:= splitiostring(ne.OutData)
		ps:= splitpstring(strs[0])
		var utxo inc.UTXO
		utxo.Owner=ps[0]
		utxo.Balance,_ = strconv.ParseFloat(ps[1], 64)
		if swdb.SW.UTXO_add(ne.Exhash+"_0",&utxo){
			fmt.Println("[INFO] Get a utxo ",ne.Exhash+"_0",utxo)
			ee.Status=1
			swdb.SW.Event_add(ee.Exhash,&ee)
			return true
		}else{
			ee.Status=0
			swdb.SW.Event_add(ee.Exhash,&ee)
			return true
		}
	}
	ee.Status=0
	swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
	return false
}


func Fund1_event_proc(ne *inc.NewEvent,height int)bool{
	ee := inc.Event{ne.Exhash,ne.EventType,ne.EventID,ne.From,ne.Sign,ne.Timestamp,ne.InputData,ne.OutData,1,height}
	if Fund1_event_check(ne){
		//删除输入中的utxo
		var input inc.UTXO_INPUTS
		json.Unmarshal([]byte(ne.InputData),&input)
		for i:=0;i<len(input.UTXOI);i++{
			swdb.SW.UTXO_del(input.UTXOI[i].Utxo)
			fmt.Println("[INFO] Del a utxo ",input.UTXOI[i].Utxo)
		}

		//循环输出，添加utxo
		var output inc.UTXO_OUTPUTS
		json.Unmarshal([]byte(ne.OutData),&output)
		for i:=0;i<len(output.UTXOO);i++{
			swdb.SW.UTXO_add(ne.Exhash+"_"+strconv.Itoa(i),&output.UTXOO[i])
			fmt.Println("[INFO] Get a utxo ",ne.Exhash+"_"+strconv.Itoa(i),output.UTXOO[i])
		}

		swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
		return true
	}
	ee.Status=0
	swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
	return false
}

func Fund2_event_proc(ne *inc.NewEvent,height int)bool{
	ee := inc.Event{ne.Exhash,ne.EventType,ne.EventID,ne.From,ne.Sign,ne.Timestamp,ne.InputData,ne.OutData,1,height}
	if Fund2_event_check(ne){
		pub,_:= hex.DecodeString(ne.Sign.Pubkey)
		addr:= GetAddress(pub)

		var output inc.ACC_PARAS
		json.Unmarshal([]byte(ne.OutData),&output)
		var out_funds float64
		for i:=0;i<len(output.ACCP);i++{ //修改目标账户状态
			acc_s := swdb.SW.Account_get(output.ACCP[i].Addr)
			acc_s.Nonce+=1
			acc_s.Bal+=output.ACCP[i].Bal
			out_funds+=output.ACCP[i].Bal
			swdb.SW.Account_set(output.ACCP[i].Addr,acc_s)
			fmt.Println("Set acc worldstatus",output.ACCP[i].Addr,acc_s)
		}
		//修改源账户状态
		acc_s := swdb.SW.Account_get(addr)
		if acc_s==nil{ //对新账户添加状态
			acc_s = &inc.AccountStatus{0,0}
		}
		acc_s.Nonce+=1
		acc_s.Bal-=out_funds
		swdb.SW.Account_set(addr,acc_s)
		fmt.Println("Set acc worldstatus",addr,acc_s)

		swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
		return true
	}

	ee.Status=0
	swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
	return false
}

func Fund3_event_proc(ne *inc.NewEvent,height int)bool{
	ee := inc.Event{ne.Exhash,ne.EventType,ne.EventID,ne.From,ne.Sign,ne.Timestamp,ne.InputData,ne.OutData,1,height}
	if Fund3_event_check(ne){
		pub,_:= hex.DecodeString(ne.Sign.Pubkey)
		addr:= GetAddress(pub)

		//循环输出，添加utxo
		var output inc.UTXO_OUTPUTS
		json.Unmarshal([]byte(ne.OutData),&output)
		var out_funds float64
		for i:=0;i<len(output.UTXOO);i++{
			swdb.SW.UTXO_add(ne.Exhash+"_"+strconv.Itoa(i),&output.UTXOO[i])
			out_funds+= output.UTXOO[i].Balance
			fmt.Println("[INFO] Get a utxo ",ne.Exhash+"_"+strconv.Itoa(i),output.UTXOO[i])
		}


		//修改源账户状态
		acc_s := swdb.SW.Account_get(addr)
		acc_s.Nonce+=1
		acc_s.Bal-=out_funds
		swdb.SW.Account_set(addr,acc_s)
		fmt.Println("Set acc worldstatus",addr,acc_s)

		swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
		return true


	}
	ee.Status=0
	swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
	return false
}

func Fund4_event_proc(ne *inc.NewEvent,height int)bool{
	ee := inc.Event{ne.Exhash,ne.EventType,ne.EventID,ne.From,ne.Sign,ne.Timestamp,ne.InputData,ne.OutData,1,height}

	if Fund4_event_check(ne){
		//删除输入中的utxo
		var input inc.UTXO_INPUTS
		json.Unmarshal([]byte(ne.InputData),&input)
		for i:=0;i<len(input.UTXOI);i++{
			swdb.SW.UTXO_del(input.UTXOI[i].Utxo)
			fmt.Println("[INFO] Del a utxo ",input.UTXOI[i].Utxo)
		}

		var output inc.ACC_PARAS
		json.Unmarshal([]byte(ne.OutData),&output)
		var out_funds float64
		for i:=0;i<len(output.ACCP);i++{ //修改目标账户状态
			acc_s := swdb.SW.Account_get(output.ACCP[i].Addr)
			if acc_s==nil{ //对新账户添加状态
				acc_s = &inc.AccountStatus{0,0}
			}
			acc_s.Nonce+=1
			acc_s.Bal+=output.ACCP[i].Bal
			out_funds+=output.ACCP[i].Bal
			swdb.SW.Account_set(output.ACCP[i].Addr,acc_s)
			fmt.Println("Set acc worldstatus",output.ACCP[i].Addr,acc_s)
		}

		swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
		return true
	}
	ee.Status=0
	swdb.SW.Event_add(ee.Exhash,&ee) //添加新事件存储
	return false
}


func handlerError(ne *inc.NewEvent){
	recover()
}