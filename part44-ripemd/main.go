package main

import (
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	hasher := ripemd160.New()
	hasher.Write([]byte("http://www.heilluo.com"))
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)
}
