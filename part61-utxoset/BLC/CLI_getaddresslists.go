package BLC

import "fmt"

//打印所有的钱包地址
func  (cli *CLI) addressLists() {
	fmt.Println("打印所有的钱包地址:")
	wallets, _  := NewWallets()
	for addresss, _ := range wallets.WalletsMap {
		fmt.Println(addresss)
	}
}
