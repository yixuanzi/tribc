package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"tribc/core"
	"tribc/inc"
	"tribc/lib"
	"tribc/swdb"
)

var current="/home/lab8/go/src/tribc/"

func main(){

	if !lib.Isgorun(os.Args[0]){
		current=lib.GetCurrent(os.Args[0])
	}
	fmt.Println(os.Args,current)

	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(1)
	var cc inc.CoreConf
	data,_ := ioutil.ReadFile(current+"conf/tri.conf")
	json.Unmarshal(data,&cc)
	core.CC=&cc
	core.CreateRPC(&cc)
	core.InitCRPC(&cc)
	swdb.InitSW(current+"data") //数据存储目录
	//go core.TestEvent()
	fmt.Println("All of thing is Over,wait it down!!!")
	caputerExit(&wg)
	wg.Wait()
	fmt.Println("Get exit signal,store and exiting...")
	swdb.SW.Close_sw()
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

