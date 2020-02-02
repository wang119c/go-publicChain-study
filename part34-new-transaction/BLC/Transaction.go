package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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

//第二种 转账时产生的 Transaction
func NewSimpleTransaction(from string, to string, amount int) *Transaction {

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	//代表消费
	bytes2, _ := hex.DecodeString("2f594c9c13526177eb2269171fc1b1eb45c1f44858421e56dc5cd0cc2a0c3566")
	txInput := &TXInput{bytes2, 0, from}
	//消费
	txInputs = append(txInputs, txInput)
	//转账
	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	//找零
	txOutput = &TXOutput{int64(10 - amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置hash值
	tx.HashTransaction()

	return tx
}
