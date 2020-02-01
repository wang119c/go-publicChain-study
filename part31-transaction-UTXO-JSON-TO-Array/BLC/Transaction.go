package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//交易数据 UTXO
type Transaction struct {
	TxHash []byte      //交易hash
	Vins   []*TXInput  //输入
	Vouts  []*TXOutput //输出
}

//1.Transaction 创建分两种情况
//第一种 创世区块创建时的 Transaction

func NewCoinbaseTransaction(address string) *Transaction {
	//代表消费
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//第二种转账时的 Transaction
