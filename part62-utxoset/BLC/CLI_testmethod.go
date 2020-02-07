package BLC

func (cli *CLI) testMethod() {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	//utxoMap := blockchain.FindUTXOMap()
	//fmt.Println(utxoMap)

	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()

}
