package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddBlock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func main() {

	//flagString := flag.String("printChain", "", "输出所有的区块")
	//flagInt := flag.Int("number", 6, "输出一个整数...")
	//flagBool := flag.Bool("open",false,"判断真假...")
	//flag.Parse()
	//fmt.Printf("%s\n", *flagString)
	//fmt.Printf("%d\n", *flagInt)
	//fmt.Printf("%v\n", *flagBool)

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "http://www.heilluo.com", "交易数据")

	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		fmt.Println(err)
		if err != nil {
			log.Panic(err)
		}
		break
	case "printChain":
		err := addBlockCmd.Parse(os.Args[2:])
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
	}

	if printChainCmd.Parsed() {
		fmt.Println("输出区块所有数据")
	}

}
