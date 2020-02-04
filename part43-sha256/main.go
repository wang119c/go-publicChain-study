package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	hasher := sha256.New()
	hasher.Write([]byte("http://www.heilluo.com"))
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)
}
