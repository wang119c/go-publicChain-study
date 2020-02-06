package BLC

import (
	"fmt"
	"os"
)

//转账
func (cli *CLI) send(from []string, to []string, amount []string) {
	if dbExists() == false {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	blockchain.MineNewBlock(from, to, amount)
	blockchain.DB.Close()
}
