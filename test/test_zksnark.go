package main

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"tribc/core"
	"tribc/lib"
)

type Commit struct {
	Addr string
	Rd *big.Int
	Mount int64
	AddrV *big.Int
}

func BigAdd(a1 *big.Int,a2 *big.Int,a3 *big.Int)*big.Int{
	temp := new(big.Int).Add(a1,a2)
	rs := temp.Add(temp,a3)
	return rs
}

func BigSub(a1 *big.Int,a2 *big.Int,a3 *big.Int) *big.Int{
	temp := new(big.Int).Sub(a1,a2)
	rs := temp.Sub(temp,a3)
	return rs
}

//计算返回a_i+r_i
func (cm *Commit) GetApR() *big.Int{
	return new(big.Int).Add(cm.AddrV,cm.Rd)
}

//计算返回R_i
func (cm *Commit) GetR() []byte{
	bm:=new(big.Int).SetInt64(cm.Mount)
	rs := BigAdd(cm.AddrV,cm.Rd,bm)
	rsbyte := rs.Bytes()
	b := make([]byte, 0, 32) // 长度为32的byte
	rsbyte = lib.PaddedAppend(32, b, rsbyte) //补齐到32
	return rsbyte

}
//计算返回 H_i
func (cm *Commit) Getsha256() []byte{
	rsbyte:=cm.GetR()
	sha:= core.ComputerHash(rsbyte)
	return sha
}

//初始化随机数Rd ，并计算账户地址字符串对应值
func (cm *Commit) SetRa() *big.Int{
	rb:=[]byte(lib.GenerateRstring(16))
	cm.Rd=new(big.Int).SetBytes(rb)

	if cm.AddrV==nil{
		md5key:=md5.Sum([]byte(cm.Addr))
		cm.AddrV=new(big.Int).SetBytes(md5key[:])
	}
	return cm.Rd
}

func main()  {
	var comm [16]byte
	fmt.Println("手动定义测试")

	R_blank:=append(comm[:],comm[:]...)
	/* R1=R2+R3+X
	r2:= new(big.Int).SetBytes([]byte("11661166abcdabcd"))
	R2:= append(comm[:],r2.Bytes()...)

	r3:= new(big.Int).SetBytes([]byte("2abcdefgzabcdefg"))
	R3:=append(comm[:],r3.Bytes()...)


	r1:= new(big.Int).SetBytes([]byte("z111z111z111z111"))
	R1:=append(comm[:],r1.Bytes()...)
	*/

	//R1+X=R2+r3
	r2:= new(big.Int).SetBytes([]byte("21661166abcdabcd"))
	R2:= append(comm[:],r2.Bytes()...)

	r3:= new(big.Int).SetBytes([]byte("2abcdefgzabcdefg"))
	R3:=append(comm[:],r3.Bytes()...)


	r1:= new(big.Int).SetBytes([]byte("3111z111z111z111"))
	R1:=append(comm[:],r1.Bytes()...)

	temp:=new(big.Int).Sub(r1,r2)
	x:=temp.Sub(temp,r3)
	X:=append(comm[:],x.Bytes()...)

	/*
	var jw [32]int
	fv:=0
	for i:=31;i>0;i--{
		fv=int((int(R1[i]) +int(R2[i])+fv)/256)
		jw[i-1]=fv
	}*/

	//fmt.Println("R1=R2+R2+X")
	fmt.Println("R1+X=R2+R2")
	fmt.Println("R1",R1)
	fmt.Println("R2",R2)
	fmt.Println("R3",R3)
	fmt.Println("X",X)

	H1:= core.ComputerHash(R1)
	H2:= core.ComputerHash(R2)
	H3:= core.ComputerHash(R3)
	H_blank:= core.ComputerHash(R_blank)

	fmt.Println("H_i=sha256(R_i)")
	fmt.Println("H1",H1)
	fmt.Println("H2",H2)
	fmt.Println("H3",H3)
	fmt.Println("================================")
	fmt.Println("Trias 业务下的测试验证样例。。。。。。")
	U1:=Commit{"abcdefg",nil,100,nil}
	U2:=Commit{"1112222",nil,80,nil}
	U3:=Commit{"abcdefg",nil,20,nil}

	U1.SetRa()
	R1=U1.GetR()
	H1=U1.Getsha256()
	ar1:=U1.GetApR()

	U2.SetRa()
	R2=U2.GetR()
	H2=U2.Getsha256()
	ar2:=U2.GetApR()

	U3.SetRa()
	R3=U3.GetR()
	H3=U3.Getsha256()
	ar3:=U3.GetApR()

	XX:=BigSub(ar1,ar2,ar3) // ar1-ar2-ar3 ==> ar_i=a_i+r_i
	b := make([]byte, 0, 32) // 长度为32的byte
	X = lib.PaddedAppend(32, b, XX.Bytes()) //补齐到32

	if XX.Sign() < 0{
		fmt.Println("R1+X=R2+R3")
	}else{
		fmt.Println("R1=R2+R3+X")
	}
	fmt.Println("R1",R1)
	fmt.Println("R2",R2)
	fmt.Println("R3",R3)
	fmt.Println("X",X)

	fmt.Println("H1",H1)
	fmt.Println("H2",H2)
	fmt.Println("H3",H3)

	fmt.Println("=========blank object==========")
	fmt.Println("R_blank",R_blank)
	fmt.Println("H_blank",H_blank)

	//fmt.Println(ar1,ar2,ar3)
}
