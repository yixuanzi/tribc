# Name: Trias Account Server(TAS)
# Data: 2018-12-27
# Version: 1.6.1
# Author: Guo Guisheng (Trias lab)

-------
# Overview
This project is mainly to provide application-level account services for trias, so that logic can also operate abstract account status through simple interface calls.

# Main Features
 - craete account
 - save account data to file
 - encryption account data in file store
 - easy and security account rpc
 - support shield address to shield transaction address
 
 
# Build
Build the project must have Golang env;
```shell
go build trias_accs.go
```

# Run
```shell
# default the server listen in 127.0.0.1:9876
 
trias_accs [ip:port]
```

# Test
We will provide the test script with python for current rpc server,use below command to test server.
```shell
# default the test script connect to 127.0.0.1:9876
# you should modify the address in code when you run rpc server in other listen address

python test/rpclient.py
```
 
# Document
Design file: doc/Trias 账号SDK设计实现.md
 
# To Do
 - support Zero—Knowledge Proof to shield transaction amount