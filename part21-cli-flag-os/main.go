package main

import (
	"fmt"
	"os"
)

func main() {

	//flagString := flag.String("printChain", "", "输出所有的区块")
	//flagInt := flag.Int("number", 6, "输出一个整数...")
	//flagBool := flag.Bool("open",false,"判断真假...")
	//flag.Parse()
	//fmt.Printf("%s\n", *flagString)
	//fmt.Printf("%d\n", *flagInt)
	//fmt.Printf("%v\n", *flagBool)

	args := os.Args
	fmt.Printf("%v\n", args)
	fmt.Printf("%v\n", args[1])

}
