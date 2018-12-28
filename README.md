# TRIBC
# Author:yixuanzi
# Data: 2018-12-28
# Email: yeying0311@26.com
# Version: 1.6.1

------
# OverView
A sample and flexible architecture blockchain with go lang.

The project have the consensus layer and application layer.

The consensus layer provide the consensus and data Synchronize Only.

The application layer provide the business logic and book store Only.

The consensus layer and application layer have to Communication with RPC calls. 


# Main Featrues
 - Flexible architecture
 - Privacy transaction

# Build & Run
```shell
//Build the project
go build consensus_main.go
go build trias_main.go
```

```shell
//runing the projcect
// init the consensus node conf
./consensus_main init

// run the consensus node 
./consensus_main run

// run the application process
./trias_main

//tips:above the command will read the conf file from $current/conf
```


# Test
Use wallet.go to test the blockchain with utxo transaction.