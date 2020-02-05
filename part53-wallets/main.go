package main

import (
	"fmt"
	"go-publicChain-study/part53-wallets/BLC"
)

func main() {
	wallets := BLC.NewWallets()
	fmt.Println(wallets.Wallets)
	wallets.CreateNewWallet()
	wallets.CreateNewWallet()
	fmt.Println(wallets.Wallets)



	//wallet := BLC.NewWallet()
	//address := wallet.GetAddress()
	//
	//fmt.Printf("addrss:%s", address)
	//isValid := wallet.IsValidForAddress(address)
	//
	//fmt.Printf("%s is %v\n", address, isValid)

}
