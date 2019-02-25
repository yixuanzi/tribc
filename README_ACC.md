# Name: Trias Account Moudel(TAM)
# Data: 2019-02-25
# Version: 1.8.0
# Author: Guo Guisheng (Trias lab)

-------
# Overview
This project is mainly to provide application-level account services for trias, so that logic can also operate abstract account status through simple interface calls.

# Main Features
 - support mult-model call from client,it's rpc,wasm and dynamic link library
 - craete account
 - save account data to file
 - encryption account data to file store
 - easy and security account rpc
 - support shield address to shield transaction address
 
# RPC Model
## Build
Build the project must have Golang env;
```shell
go build trias_accs.go
```

## Run
```shell
# default the server listen in 127.0.0.1:9876
 
trias_accs [ip:port]
```

## Test
We will provide the test script with python for current rpc server,use below command to test server.
```shell
# default the test script connect to 127.0.0.1:9876
# you should modify the address in code when you run rpc server in other listen address

python test/rpclient.py
```

# WASM Model
this model is support the browser to call the exist sdk lib
## Build
```shell
#require go >= 1.11
./buildwasm.sh triacc_wasm
```
 
## Test
```shell
1. go run http.go (To start a http server)
2. use chrome open localhost:8080/wasm_exec.html,click the button to test the triacc function call
```

# Dynamic link library
this model is support the local client application,android device application。
## Build
```shell
NONE
``` 
 
# Document
Design file: doc/Trias 账号SDK设计实现.md
 
# To Do
 - support Zero—Knowledge Proof to shield transaction amount
