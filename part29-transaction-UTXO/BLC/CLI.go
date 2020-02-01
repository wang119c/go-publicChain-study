package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address DATA -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(txs []*Transaction) {
	if dbExists() == false {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	blockchain.AddBlockToBlockchain(txs)

	blockchain.DB.Close()
}

func (cli *CLI) printchain() {

	if dbExists() == false {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	blockchain.Printchain()

	blockchain.DB.Close()
}

func (cli *CLI) createGenesisBlockchain(address string) {
	CreateBlockchainWithGenesisBlock(address)
}

func (cli *CLI) Run() {
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "http://www.heilluo.com", "交易数据")
	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块地址")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break
	default:
		printUsage()
		os.Exit(1)
		break
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagAddBlockData)
		cli.addBlock([]*Transaction{})
	}

	if printChainCmd.Parsed() {
		//fmt.Println("输出区块所有数据")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空...")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)

		////创世区块
		//blockchain := BLC.CreateBlockchainWithGenesisBlock()
		//defer blockchain.DB.Close()
	}

}
