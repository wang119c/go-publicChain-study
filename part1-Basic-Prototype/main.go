package main

import (
	"fmt"
	"go-publicChain-study/part1-Basic-Prototype/BLC"
)

func main() {
	block := BLC.CreateGenesisBlock(	"Genenis Block")
	fmt.Println(block)

}
