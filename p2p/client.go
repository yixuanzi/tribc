package p2p

import (
	"bytes"
	"fmt"
	"net"
	"time"
)


type Dstp2p struct{
	Addr string
	Conn *net.TCPConn
	Msg []byte
	IP string
}

func (dst *Dstp2p)ConnServer()bool{
	//var buf [1024*1024*1]byte
	defer connFail(dst)
	time.Sleep(time.Duration(3)*time.Second)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", dst.Addr)
	checkErr(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//defer conn.Close()
	checkErr(err)
	dst.Conn=conn
	rAddr := conn.RemoteAddr()
	_, err = conn.Write([]byte("TRIHello server!TRI"))
	checkErr(err)
	fmt.Println("Conn Server succ",rAddr.String())
	/*

	for {
		n, err = conn.Read(buf[0:])
		checkErr(err)
		fmt.Println("Reply from server ", rAddr.String(), string(buf[0:n]))
		handleClient(buf)

	}
	os.Exit(0)
	*/
	return true
}

//广播数据，根据事件类型构建广播数据结构
func (dst *Dstp2p)Sendmsg(mess []byte,msgtype uint8)bool{
	defer sendFail(dst)
	if dst.Conn==nil{
		fmt.Println("[WARN] Sendmsg error,the for the server tcp connect is unvaild")
		return false
	}
	var bcdata bytes.Buffer
	bcdata.Write([]byte("TRI"))
	bcdata.WriteByte(msgtype)
	bcdata.Write(mess)
	bcdata.Write([]byte("TRI"))
	_, err := dst.Conn.Write(bcdata.Bytes())
	checkErr(err)
	return true
}

func connFail(dst *Dstp2p){
	if err := recover(); err != nil {
		fmt.Println("Conn Server Fail ",dst.Addr,err)
	}
}

func sendFail(dst *Dstp2p){
	if err := recover(); err != nil {
		fmt.Println("[Error] Sendmsg fail ",dst.Addr,err)
	}
}