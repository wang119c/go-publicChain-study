package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//1，有一个办法 ， 功能:
//遍历整个数据库，读取所有的未花费的UTXO , 然后将所有的UTXO存储到数据库
//去遍历数据库时

const utxoTableName = "utxoTableName"

type UTXOSet struct {
	Blockchain *Blockchain
}

//重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet() {
	err := utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))

			if err != nil {
				log.Panic(err)
			}

		}

		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			txOutputsMap := utxoSet.Blockchain.FindUTXOMap()

			for keyHash, outs := range txOutputsMap {

				txHash, _ := hex.DecodeString(keyHash)
				b.Put(txHash, outs.Serialize())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (utxoSet *UTXOSet) findUTXOForAddress(address string) []*UTXO {

	var utxos []*UTXO

	utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		//游标
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("key = %s ,value = %v\n", k, v)
			txOutputs := DeserializeTXOutPuts(v)
			for _, utxo := range txOutputs.UTXOS {
				if utxo.Output.UnlockScriptPubKeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}
		return nil
	})
 	return utxos
}

//查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int64 {

	UTXOS := utxoSet.findUTXOForAddress(address)
	var amount int64

	for _, utxo := range UTXOS {
		amount += utxo.Output.Value
	}

	return amount
}
