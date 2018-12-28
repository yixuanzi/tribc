package main

import (
	"fmt"
	"tribc/inc"
	"tribc/swdb"
)

func main() {
	fmt.Println("this is test with swdb")
	sw := new(swdb.StatusWorld)

	u:= inc.UTXO{"abc",1.1}

	fmt.Println(u)

	sw.Init_sw("/tmp/gkvdb")
	sw.UTXO_add("a",&u)
	sw.UTXO_add("b",&u)

	fmt.Println(sw.UTXO_get("a"))
	fmt.Println(sw.UTXO_get("b"))
	sw.Close_sw()
}
