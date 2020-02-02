package main

import (
	"go-publicChain-study/part38-new-transaction/BLC"
)

func main() {

	cli := BLC.CLI{}
	cli.Run()

	//./main createblockchain -address "helloo"
	// ./main printchain
	// ./main send -from '["liyue","zhangxqing"]' -to '["juncheng","xiaoyong"]' -amount '["2","3"]'

	////新区块
	//blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang")
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang")
	//blockchain.AddBlockToBlockchain("Send 200RMB To zhangqiang")
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang")
	//
	//blockchain.Printchain()

}
