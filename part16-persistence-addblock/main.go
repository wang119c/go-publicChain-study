package main

import (
	"go-publicChain-study/part16-persistence-addblock/BLC"
)

func main() {
	//创世区块
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()
	//新区块
	blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang")
	blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang")
	blockchain.AddBlockToBlockchain("Send 200RMB To zhangqiang")
	blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang")




}
