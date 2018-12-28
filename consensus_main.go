package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"tribc/p2p"
	"tribc/lib"
)

var current="/home/lab8/go/src/tribc"

func main(){
	if !lib.Isgorun(os.Args[0]){
		current=lib.GetCurrent(os.Args[0])
	}
	fmt.Println(os.Args,current)

	if len(os.Args)<2{
		fmt.Println("Please input vaild parameter!")
		os.Exit(1)
	}
	if os.Args[1]=="init"{
		p2p.Init(current+"/conf/consensus.conf")
		os.Exit(0)

	}else if os.Args[1]=="run"{

		runtime.GOMAXPROCS(1)
		var wg sync.WaitGroup
		wg.Add(1)

		node:= p2p.Load(current+"/conf/consensus.conf",current+"/conf/p2p.conf")
		fmt.Println("load conf file succ",node)

		fmt.Println("Create P2PNet...")
		pnet:= p2p.CreateNet(node)

		fmt.Println("Start P2P listen Succ...")
		pnet.CreateServer()

		fmt.Println("Start Connect P2P node ...")
		pnet.ConnServer()

		//time.Sleep(time.Duration(1)*time.Second)
		//pnet.Broadcast([]byte("this is a test message"))
		fmt.Println("Create RPC Server ...")
		p2p.CreateRPC(node)

		fmt.Println("Connect RRPC Server ...")
		p2p.Initcp2prpc("127.0.0.1:8898")

		go pnet.Generateblock()

		fmt.Println("All of thing is Over,wait it down!!!")
		caputerExit(&wg)
		//fmt.Println(time.Now().Unix())
		wg.Wait()
		//you can do something for exit action

	}else{
		fmt.Println("Please input vaild parameter!")
		os.Exit(1)
	}

	//time.Sleep(time.Duration(3600)*time.Second)
}


func caputerExit(wg *sync.WaitGroup){
	stopChan := make(chan struct{},1)
	signalChan := make(chan os.Signal,1)
	go func(){
		<-signalChan
		stopChan<- struct{}{}
		wg.Done()
	}()
	signal.Notify(signalChan,syscall.SIGINT,syscall.SIGTERM)
}
