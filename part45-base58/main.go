package main

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

func main() {
	bytes := []byte("http://www.liyuechun.org")
	fmt.Println(base58.Encode(bytes))

	fmt.Println(base58.Decode(base58.Encode(bytes)))

}
