# Trias Account Moudel(TAM)
[![Tag version](https://img.shields.io/badge/Tag-1.7.4-blue.svg)]()
[![Go doc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/trias-lab/tribc)
[![Go version](https://img.shields.io/badge/go-1.11-blue.svg)](www.golang.org)
[![Lincese](https://img.shields.io/badge/Lincese-GPL3.0-blue.svg)](www.golang.org)

# Overview
This project is mainly to provide application-level account services for trias, so that logic can also operate abstract account status through simple interface calls.

|branch|status|test|
|-------|--------|----|
|master| latest code and function in this|
|SameA| TriasV1 working at current branch  |

# Main Features
 - support mult-model call from client,it's rpc,wasm and dynamic link library

 - craete account

 - save account data to file

 - encryption account data to file store

 - easy and security account rpc

 - support shield address to shield transaction address

# Requirements

| Requirement | Notes           |
| ----------- | --------------- |
| Go          | 1.11 or highter |
| Python      | 3.6             |




# RPC Model
## Build
Build the project must have Golang env;
```shell
go build trias_accs.go
```

## Start
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
this model is support the browser to call the exist function sdk
## Build
```shell
#require go >= 1.11
./buildwasm.sh triacc_wasm
```

## Test
```shell
go run http.go (To start a http server)

use chrome open localhost:8080/wasm_exec.html,click the button to test the triacc function call
```

# To Do
 - support Zero—Knowledge Proof to shield transaction amount


# A Note on Production Readiness

While Trias is being used in production in private, permissioned
environments, we are still working actively to harden and audit it in preparation
for use in public blockchains.
We are also still making breaking changes to the protocol and the APIs.
Thus, we tag the releases as *alpha software*.

In any case, if you intend to run Trias in production,
please [contact us](mailto:contact@trias.one) and [join the chat](https://www.trias.one).

# Security

To report a security vulnerability,  [bug report](mailto:contact@trias.one)

# Documentation
Interface design document  at [Interface design document with TAM.md](doc/Interface%20design%20document%20with%20TAM.md)

Privacy Transaction Design Principle Document at [Privacy transaction with TAM.md](Privacy%20transaction%20with%20TAM.md)

Complete documentation can be found on the [website](https://github.com/trias-lab/Documentation).



# Contributing
All code contributions and document maintenance are temporarily responsible for TriasLab

Trias are now developing at a high speed and we are looking forward to working with quality partners who are interested in Trias. If you want to join.Please contact us:
- [Telegram](https://t.me/triaslab)
- [Medium](https://medium.com/@Triaslab)
- [BiYong](https://0.plus/#/triaslab)
- [Twitter](https://twitter.com/triaslab)
- [Gitbub](https://github.com/trias-lab/Documentation)
- [Reddit](https://www.reddit.com/r/Trias_Lab)
- [More](https://www.trias.one/)
- [Email](mailto:contact@trias.one)


## Upgrades

Trias is responsible for the code and documentation upgrades for all Trias modules.
In an effort to avoid accumulating technical debt prior to Beta,
we do not guarantee that data breaking changes (ie. bumps in the MINOR version)
will work with existing Trias blockchains. In these cases you will
have to start a new blockchain, or write something custom to get the old data into the new chain.

# Resources
## Research

* [The latest paper](https://www.contact@trias.one/attachment/Trias-whitepaper%20attachments.zip)
* [Project process](https://trias.one/updates/project)
* [Original Whitepaper](https://trias.one/whitepaper)
* [News room](https://trias.one/updates/recent)
