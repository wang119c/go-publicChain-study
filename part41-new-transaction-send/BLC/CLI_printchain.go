package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printchain() {

	if dbExists() == false {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	blockchain.Printchain()

	blockchain.DB.Close()
}
