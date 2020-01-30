package main

import (
	"fmt"
	"go-publicChain-study/part14-proof-of-work/BLC"
)

func main() {
	////创世区块
	//blockchain := BLC.CreateBlockchainWithGenesisBlock()
	////新区块
	//blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 200RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//
	//fmt.Println(blockchain.Blocks)

	block := BLC.NewBlock("Test",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	fmt.Printf("%d\n",block.Nonce)
	fmt.Printf("%x\n",block.Hash)

	proofOfWork :=  BLC.NewProofOfWork(block)
	fmt.Printf("%v",proofOfWork.IsValid())





}
