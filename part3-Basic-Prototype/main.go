package main

import (
	"fmt"
	"go-publicChain-study/part3-Basic-Prototype/BLC"
)

func main() {
	genesisBlockchain := BLC.CreateBlockchainWithGenesisBlock()
	fmt.Println(genesisBlockchain)
}
