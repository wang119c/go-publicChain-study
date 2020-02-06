package BLC

import "fmt"

//先用他去查询余额
func (cli *CLI) getBalance(address string) {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)
	fmt.Printf("%s 一共有 %d 个 token \n", address, amount)
}
