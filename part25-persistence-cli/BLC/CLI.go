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
	fmt.Println("\tcreateblockchain -data DATA -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	AddBlockToBlockchain(data)
}

func (cli *CLI) printchain() {
	cli.Blockchain.Printchain()
}

func (cli *CLI) createGenesisBlockchain(data string) {
	CreateBlockchainWithGenesisBlock(data)
}

func (cli *CLI) Run() {
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "http://www.heilluo.com", "交易数据")
	flagCreateBlockchainWithData := createBlockchainCmd.String("data", "Genesis Data...", "创建创世区块交易")

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
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		//fmt.Println("输出区块所有数据")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithData == "" {
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithData)

		////创世区块
		//blockchain := BLC.CreateBlockchainWithGenesisBlock()
		//defer blockchain.DB.Close()
	}

}
