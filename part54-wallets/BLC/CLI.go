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
	fmt.Println("\tcreatewallet  -- 创建钱包")
	fmt.Println("\tcreateblockchain -address DATA -- 创建创世区块")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tgetbalance -address -- 输出区块信息")
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

func (cli *CLI) Run() {
	isValidArgs()

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块地址")
	getbalanceWithAddress := getbalanceCmd.String("address", "", "查询某一个账号的余额")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
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
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		break

	default:
		printUsage()
		os.Exit(1)
		break
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagFrom)
		//fmt.Println(*flagTo)
		//fmt.Println(*flagAmount)
		//
		//fmt.Println(JSONToArray(*flagFrom))
		//fmt.Println(JSONToArray(*flagTo))
		//fmt.Println(JSONToArray(*flagAmount))

		//cli.addBlock([]*Transaction{})

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)

		for index, fromAddress := range from {
			if IsValidForAddress([]byte(fromAddress)) == false || IsValidForAddress([]byte(to[index])) == false {
				fmt.Println("地址不合法......")
				os.Exit(1)
			}
		}

		cli.send(from, to, amount)
	}

	if printChainCmd.Parsed() {
		//fmt.Println("输出区块所有数据")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {
		if IsValidForAddress([]byte(*flagCreateBlockchainWithAddress)) == false {
			fmt.Println("地址不合法...")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)

		////创世区块
		//blockchain := BLC.CreateBlockchainWithGenesisBlock()
		//defer blockchain.DB.Close()
	}

	if getbalanceCmd.Parsed() {
		if IsValidForAddress([]byte(*getbalanceWithAddress)) == false  {
			fmt.Println("地址不合法...")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAddress)
	}

	if createWalletCmd.Parsed() {
		//创建钱包
		cli.createWallet()
	}

}
