package p2p

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}

type Msg struct{
	Addr string
	Mess []byte
}

/*
针对大块数据的整流处理？
type MsgBuf struct {
	Mess bytes.Buffer
	time time.Time
}

var msgmap = make(map[string] *MsgBuf)
*/

var flag =[]byte("TRI")

func (pn *P2PNet) CreateServerRouter() bool{
	defer handler("CreateServerRouter")
	//wg.Done()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", pn.node.Lis)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	fmt.Println("CreateServer succ",pn.node.Lis)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go pn.handleServer(conn)
		//time.Sleep(time.Duration(3)*time.Second)
	}

	return true
}

func (pn *P2PNet)handleServer(conn net.Conn) {
	//defer handler("handleServer")

	var buf [1024*1024*1] byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()

		if n>6{
			if bytes.Equal(buf[:3],flag) || bytes.Equal(buf[n-3:],flag){ //完整消息，直接处理
				msg:= Msg{strings.Split(rAddr.String(),":")[0],buf[3:n-3]}
				pn.HandleMsg(&msg)
			}
			//.....，其他一次无法读取所有数据的分块拼接处理
		}

		//fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		/*
		_, err2 := conn.Write([]byte("Welcome client!"))
		if err2 != nil {
			return
		}*/
	}
	fmt.Println("handleServer over")
}

func handler(m string){
	if err := recover(); err != nil {
		fmt.Println("[Error]"+m ,err)
	}
}