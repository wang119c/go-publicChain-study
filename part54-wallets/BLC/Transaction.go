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

//判断当前的交易是否是coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].Txhash) == 0 && tx.Vins[0].Vout == -1
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
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {
	money, spendableUTXODic := blockchain.FindSpentableUTXOS(from, amount, txs)
	var txInputs []*TXInput
	var txOutputs []*TXOutput
	for txHash, indexArray := range spendableUTXODic {
		for _, index := range indexArray {
			//代表消费
			txHashBytes, _ := hex.DecodeString(txHash)
			txInput := &TXInput{txHashBytes, index, from}
			//消费
			txInputs = append(txInputs, txInput)
		}
	}

	//转账
	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	//找零
	txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置hash值
	tx.HashTransaction()

	return tx
}
