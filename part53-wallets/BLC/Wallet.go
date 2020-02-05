package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumlen = 4

type Wallet struct {
	//1.私钥
	PrivateKey ecdsa.PrivateKey
	//2.公钥
	PublicKey []byte
}

//创建一个钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	fmt.Println(privateKey, publicKey)

	return &Wallet{privateKey, publicKey}
}

//通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

//判断地址是不是有效
func (w *Wallet) IsValidForAddress(address []byte) bool {
	version_public_checksumBytes := Base58Decode(address)
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumlen:]
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumlen]
	checkBytes := CheckSum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

//生成钱包地址
func (w *Wallet) GetAddress() []byte {
	ripemd160Hash := w.Ripemd160Hash(w.PublicKey)
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)
	checkSumBytes := CheckSum(version_ripemd160Hash)
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	return Base58Encode(bytes)
}

func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	//256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	//160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	return ripemd160.Sum(nil)
}

//checksum 为一个公钥生成 checksum
func CheckSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumlen]
}
