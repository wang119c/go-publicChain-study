package BLC

import "bytes"

type TXInput struct {
	//1.交易id
	Txhash []byte
	//2.存储txoutput 在vout 里面的索引
	Vout int
	//数字签名
	Signature []byte
	//公钥
	Pubkey []byte
}

//判断当前的消费是谁的钱
func (txInput *TXInput) UnlockRipemd160Hash(ripemd160Hash []byte) bool {
	pubkey :=   Ripemd160Hash(txInput.Pubkey)
	return bytes.Compare(pubkey , ripemd160Hash) == 0
}
