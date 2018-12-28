package inc

type CoreConf struct{
	LRPC string
	RRPC string
}


// 请求参数结构体
type NERequest struct {
	NE NewEvent
	Src string
}

type TBRequest struct {
	TB TriBlock
	Src string
}

type CMDRequest struct {
	Cmd string
}

// rpc响应结构体
type RPCResponse struct {
	RS string
}

type HeightResponse struct {
	Height int
	Hash string
}

type UTXOResponse struct {
	UOBJ	UTXO
}

//--------------------------

type TSign struct {
	Pubkey string
	Sigadata string
}

type UTXO struct {
	Owner    string
	Balance  float64
}

type Sblock struct {
	Hash string
	ParentHash string
	Height int
	Miner string
	Sign TSign
	StatusRoot string
	EventRoot string
	Timestamp int64
	Eventlisthash []string
}




/*
0:funds 0:createcoin 1:utxo->utxo 2:acc->acc 3:acc->utxo 4:utxo->acc
1:account
2:contract

 */
type NewEvent struct {
	Exhash string
	EventType uint8 // 0:funds 1:account 2:contract
	EventID uint8
	From string
	Sign TSign
	Timestamp int64
	InputData string
	OutData string
}

type Event struct{
	Exhash string
	EventType uint8 // 0:funds 1:account 2:contract
	EventID uint8
	From string
	Sign TSign
	Timestamp int64
	InputData string
	OutData string
	Status uint8 //0  1
	Height int
}

type TriBlock struct {
	Hash string
	ParentHash string
	Height int
	Miner string
	Sign TSign
	StatusRoot string
	EventRoot string
	Timestamp int64
	Eventlist []NewEvent
}

type AccountStatus struct {
	Nonce int
	Bal	float64
}



//事件中的输入输出参数结构
type UTXO_INPUT struct {
	Utxo string
	Sign TSign
}



type UTXO_INPUTS struct {
	UTXOI []UTXO_INPUT
}

type UTXO_OUTPUTS struct {
	UTXOO []UTXO
}

type ACC_PARA struct {
	Addr string
	Bal float64
}

type ACC_PARAS struct {
	ACCP []ACC_PARA
}

