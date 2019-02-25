package main

import (
	"encoding/hex"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime/debug"
	"tribc/core"
	"tribc/inc"
)


type AAccount struct {
	Acc *core.Account
	Pass string
}

const version  = "1.8.0"
var AccM map[string]AAccount = map[string]AAccount{}

func main() {
	addr:="127.0.0.1:9876"
	if len(os.Args)>1{
		addr=os.Args[1]
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Start RPC Server fail,Please check the input parameters!")
		return
	}
	log.Println("Start RPC Server with trias account in ",addr,"  Version:",version)
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("AccRPC", new(AccRPC)); err != nil {
		return
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
			continue
		}
		go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}


type AccRPC struct {
}


//=================================================================
type Json struct {
	Name string `json:name`
	Age  int    `json:age`
}
func (self *AccRPC) JsonTest(args Json, result *Json) error {
	log.Println(args)
	*result = Json{args.Name, args.Age}
	return nil
}

func (self *AccRPC) Test(args string, result *string) error {
	log.Println(args)
	*result="OK"
	return nil
}
//===========================
type CAargs struct {
	Path string `json:path`
	Pass string `json:pass`
}
//创建并导出账号
func (self *AccRPC) CreateAcc(args CAargs, result * string)error{
	log.Println(args)
	defer handlerError("CreateAcc")

	acc:=*core.CreateAccount()
	addr := core.GetAddress(acc.GkeyA.GetPubKey(),acc.GkeyB.GetPubKey())

	aacc := AAccount{&acc,args.Pass}
	if _,ok := AccM[addr]; !ok{
		AccM[addr]=aacc
		if core.Save2file(&acc,args.Path,[]byte(args.Pass)){
			*result=addr
			return nil
		}
	}
	*result="Fail"
	return nil
}

func (self *AccRPC) GetAcclist(args string,result *[]string)error{
	log.Println(args)
	defer handlerError("GetAcclist")
	var addrlist []string
	for addr := range AccM{
		addrlist=append(addrlist,addr)
	}
	*result=addrlist
	return nil
}

type IAargs struct {
	Path string `json:path`
	Pass string `json:pass`
}
func (self *AccRPC) ImportAcc(args IAargs,result *string) error{
	log.Println(args)
	defer handlerError("ImportAcc")
	acc:= core.Load4file(args.Path,[]byte(args.Pass))
	if acc!=nil{
		addr := core.GetAddress(acc.GkeyA.GetPubKey(),acc.GkeyB.GetPubKey())
		aacc := AAccount{acc,args.Pass}
		if _,ok := AccM[addr];ok{
			*result="The Acc is exist!"
			return nil
		}
		*result=addr
		AccM[addr]=aacc
		return nil
	}
	*result="Fail"
	return nil
}

type EAargs struct {
	Addr string `json:addr`
	Path string `json:path`
	Pass1 string `json:pass1`
	Pass2 string `json:pass2`
}
func (self *AccRPC) ExportAcc (args EAargs,result *string) error{
	log.Println(args)
	defer handlerError("ExportAcc")
	if aacc,ok := AccM[args.Addr];ok{
		if aacc.Pass==args.Pass1{
			path:=args.Path+"/"+args.Addr+".json" //拼接形成导出账号文件路径
			pass:=args.Pass2 //支持传递参数修改导出账号文件密码
			if len(pass)==0{
				pass=args.Pass1
			}
			if core.Save2file(aacc.Acc,path,[]byte(pass)){
				*result=path
				return nil
			}
		}

	}
	*result="Fail"
	return nil
}

//===================================
type Sargs struct {
	Addr string `json:addr`
	Hash string `json:hash`
	Pass string `json:pass`
}
func (self *AccRPC)Sign(args Sargs,result *inc.TSign) error{
	log.Println(args)
	defer handlerError("Sign")
	if aacc,ok:= AccM[args.Addr];ok{
		if args.Pass==aacc.Pass{
			signdata,_:= core.Sign(aacc.Acc.GkeyB.PrivateKey,[]byte(args.Hash))
			tsign:=inc.TSign{hex.EncodeToString(aacc.Acc.GkeyB.GetPubKey()),signdata}
			*result=tsign
			return nil
		}
	}
	*result=inc.TSign{"",""}
	return nil
}

type Vargs struct{
	Pubkey string `json:pubkey`
	Hash string `json:hash`
	Stext string `json:stext`
}

func (self *AccRPC)Verify(args Vargs,result *string)error  {
	log.Println(args)
	defer handlerError("Verify")
	pub_b,_:= hex.DecodeString(args.Pubkey)
	pubk:= core.Pub2pubKey(pub_b) //当前公钥数据未解码
	f,err:=core.Verify([]byte(args.Hash),args.Stext,pubk)
	if f{
		*result="OK"
		return nil
	}
	*result="Fail"
	return err
}
//==========================================================
type RepCSA struct {
	Shieldaddr string //`json:shieldaddr`
	ShieldpKey string //`json:shieldpkey`
}
func (self *AccRPC)CreateShieldAddr(args string,result *RepCSA) error{
	log.Println(args)
	defer handlerError("CreateShieldAddr")
	shieldaddr, shieldpKey := core.CreateShieldAddr(args)
	if shieldpKey==nil{
		*result=RepCSA{"",""}
	}else{
		*result=RepCSA{hex.EncodeToString(shieldaddr),hex.EncodeToString(shieldpKey)}
	}
	return nil
}

type VSargs struct {
	Addr string `json:addr`
	Shieldaddr string `json:shieldaddr`
	ShieldpKey string `json:shieldpkey`
}
func (self *AccRPC)Verify_shield(args VSargs,result *string)error{
	log.Println(args)
	defer handlerError("Verify_shield")
	if aacc,ok:= AccM[args.Addr];ok {
		acc:=aacc.Acc
		saddr,_:=hex.DecodeString(args.Shieldaddr)
		spkey,_:=hex.DecodeString(args.ShieldpKey)
		if core.Verify_shield(acc,saddr,spkey) {
			*result="OK"
			return nil
		}
	}
	*result="Fail"
	return nil
}

type SSargs struct {
	Addr string `json:addr`
	Pass string `json:pass`
	Hash string `json:hash`
	ShieldpKey string `json:shieldpkey`

}
func (self *AccRPC)Shield_Sign(args SSargs,result *inc.TSign)error{
	log.Println(args)
	defer handlerError("Shield_Sign")
	if aacc,ok:= AccM[args.Addr];ok {
		if aacc.Pass!=args.Pass{
			*result=inc.TSign{"",""}
			return nil
		}
		acc := aacc.Acc
		spkey,_:=hex.DecodeString(args.ShieldpKey)
		priv := core.Getprivkey(acc,spkey)
		signdata,_ := core.Sign(priv,[]byte(args.Hash))
		pubKey := append(priv.PublicKey.X.Bytes(), priv.Y.Bytes()...) // []bytes type
		tsign:=inc.TSign{hex.EncodeToString(pubKey),signdata}
		*result=tsign
		return nil
	}
	*result=inc.TSign{"",""}
	return nil
}

//根据地址返回验证公钥B，用于在区块链中检验当前签名是否是当前地址的签名（签名检查分两步：签名有效性检查，当前签名用户地址和当前资产地址检查;使用账户地址返回当前账户公钥B，和当前签名的公钥进行匹配，完成所有权验证）
func (self *AccRPC)GetACCPubkB(args string,result *string)error{
	log.Println(args)
	defer handlerError("GetACCPubkB")
	_,pubB := core.GetPubk4Addr(args)
	pubB_b := append(pubB.X.Bytes(), pubB.Y.Bytes()...)
	*result = hex.EncodeToString(pubB_b)
	return nil
}

//根据公钥返回地址(仅用于隐私交易下隐藏地址交易转化），用于在区块链中检验当前签名是否是当前地址的签名
func (self *AccRPC)Shield_Pubkey2Addr(args string,result *string)error{
	log.Println(args)
	defer handlerError("Shield_Pubkey2Addr")
	*result=core.Shield_Pubkey2Addr(args)
	return nil
}

func handlerError(name string) {
	if p := recover(); p != nil {
		log.Printf("[Error] %s call error: %v",name,p)
		debug.PrintStack()
	}
}