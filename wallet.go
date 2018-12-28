package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"tribc/core"
	"tribc/inc"
	"tribc/p2p"
)

func main(){
	//===========utxo->utxo
	fmt.Println(os.Args)
	if os.Args[1]=="0"{
	i_utxo:= os.Args[2] //需要使用的utxo
	addr:= os.Args[3]   //目标地址
	bal:= os.Args[4]   //转账金额
	path:=os.Args[5]    //账号文件
	passwd:=os.Args[6]  //账号文件密码

	acc := core.Load4file(path,[]byte(passwd))
	p2p.Initcp2prpc("127.0.0.1:8898")

	utxo := p2p.CPR.GetUTXO(i_utxo)

	var input inc.UTXO_INPUTS
	_,tsign := core.ComputerHashSign([]byte(i_utxo),acc.GkeyA)
	input.UTXOI=append(input.UTXOI,inc.UTXO_INPUT{i_utxo,*tsign})
	data_i,_:=json.Marshal(input)

	var output inc.UTXO_OUTPUTS
	f, _ := strconv.ParseFloat(bal, 64)
	output.UTXOO=append(output.UTXOO,inc.UTXO{addr,f})
	if utxo!=nil && utxo.Balance-f>0 { //找零
		output.UTXOO=append(output.UTXOO,inc.UTXO{utxo.Owner,utxo.Balance-f})
	}
	data_o,_:=json.Marshal(output)

	ne := inc.NewEvent{"",0,1,"",inc.TSign{"",""},time.Now().Unix(),string(data_i),string(data_o)}

	data,_:=json.Marshal(ne)
	hash,_:= core.ComputerHashSign(data,nil)
	ne.Exhash=hash
	fmt.Println("Generate a utxo tansform event ",ne)

	//测试验证签名
	/*
	pub,_:=hex.DecodeString(tsign.Pubkey)
	hash_b := core.ComputerHash([]byte(i_utxo))
	pubkey := core.Pub2pubKey(pub)
	ff,_ := core.Verify(hash_b,tsign.Sigadata,pubkey)
	fmt.Println(ff)
	*/

	p2p.CPR.PutEvent(&ne,"127.0.0.1")
	}else if os.Args[1]=="1"{  //acc-> acc
		addr:= os.Args[2]   //目标地址
		bal:= os.Args[3]   //转账金额
		path:=os.Args[4]    //账号文件
		passwd:=os.Args[5]  //账号文件密码
		acc := core.Load4file(path,[]byte(passwd))
		p2p.Initcp2prpc("127.0.0.1:8898")

		var output inc.ACC_PARAS
		f, _ := strconv.ParseFloat(bal, 64)
		output.ACCP=append(output.ACCP,inc.ACC_PARA{addr,f})

		//from:=core.GetAddress(pub)
		data,_:=json.Marshal(output)

		ne := inc.NewEvent{"",0,2,"",inc.TSign{"",""},time.Now().Unix(),"",string(data)}

		data,_=json.Marshal(ne)
		hash,sign := core.ComputerHashSign(data,acc.GkeyA)

		ne.Exhash=hash
		ne.Sign=*sign
		fmt.Println("Generate a acc tansform event ",ne)
		p2p.CPR.PutEvent(&ne,"127.0.0.1")

	}else if os.Args[2]=="2"{ //acc->utxo
		addr:= os.Args[2]   //目标地址
		bal:= os.Args[3]   //转账金额
		path:=os.Args[4]    //账号文件
		passwd:=os.Args[5]  //账号文件密码
		acc := core.Load4file(path,[]byte(passwd))
		p2p.Initcp2prpc("127.0.0.1:8898")

		var output inc.UTXO_OUTPUTS
		f, _ := strconv.ParseFloat(bal, 64)
		output.UTXOO=append(output.UTXOO,inc.UTXO{addr,f})

		data,_:=json.Marshal(output)
		ne := inc.NewEvent{"",0,3,"",inc.TSign{"",""},time.Now().Unix(),"",string(data)}

		data,_=json.Marshal(ne)
		hash,sign := core.ComputerHashSign(data,acc.GkeyA)

		ne.Exhash=hash
		ne.Sign=*sign
		fmt.Println("Generate a acc->utxo tansform event ",ne)
		p2p.CPR.PutEvent(&ne,"127.0.0.1")

	}else if os.Args[3]=="3"{ //utxo->acc
		i_utxo:= os.Args[2] //需要使用的utxo
		addr:= os.Args[3]   //目标地址
		bal:= os.Args[4]   //转账金额
		path:=os.Args[5]    //账号文件
		passwd:=os.Args[6]  //账号文件密码

		acc := core.Load4file(path,[]byte(passwd))
		p2p.Initcp2prpc("127.0.0.1:8898")


		_,tsign := core.ComputerHashSign([]byte(i_utxo),acc.GkeyA)

		var input inc.UTXO_INPUTS
		input.UTXOI=append(input.UTXOI,inc.UTXO_INPUT{i_utxo,*tsign})
		data_i,_:=json.Marshal(input)

		var output inc.ACC_PARAS
		f, _ := strconv.ParseFloat(bal, 64)
		output.ACCP=append(output.ACCP,inc.ACC_PARA{addr,f})
		data_o,_:=json.Marshal(output)


		ne := inc.NewEvent{"",0,4,"",inc.TSign{"",""},time.Now().Unix(),string(data_i),string(data_o)}

		data,_:=json.Marshal(ne)
		hash,_:= core.ComputerHashSign(data,nil)
		ne.Exhash=hash
		fmt.Println("Generate a utxo->acc tansform event ",ne)
		p2p.CPR.PutEvent(&ne,"127.0.0.1")

	}else {
		fmt.Println("Please input vaild parameters!")
	}
}