package main

import (
	"fmt"
	"go-publicChain-study/part50-wallet/BLC"
)

func main() {
	wallet := BLC.NewWallet()
	address := wallet.GetAddress()

	fmt.Printf("addrss:%s", address)
	isValid := wallet.IsValidForAddress(address)

	fmt.Printf("%s is %v\n", address, isValid)

}
