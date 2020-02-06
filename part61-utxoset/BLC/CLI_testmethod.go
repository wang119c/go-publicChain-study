package BLC

import "fmt"

func (cli *CLI) testMethod() {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	utxoMap := blockchain.FindUTXOMap()
	fmt.Println(utxoMap)
}
