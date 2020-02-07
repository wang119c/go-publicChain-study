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

func (utxoSet *UTXOSet) FindUnPackageSpentableUTXOS(from string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(from))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]

				if in.UnlockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.Txhash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	work1:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddress(from) {
				if spentTXOutputs != nil {
					if len(spentTXOutputs) != 0 {

						var isSpentUTXO bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								if index == i && txHash == hex.EncodeToString(tx.TxHash) {
									isSpentUTXO = true
									continue work1
								}
							}
						}

						if isSpentUTXO == false {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					} else {
						utxo := &UTXO{tx.TxHash, index, out}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			}
		}

	}

	for _, tx := range txs {
	work3:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddress(from) {

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {
							var isUnSpentUTXO bool

							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue work3
								}
							}

							if isUnSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}

						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}
		}
	}
	return unUTXOs
}

//返回要凑多少钱  ，返回对应的 txoutput  hash 及 index

func (utxoSet *UTXOSet) FindSpentableUTXOS(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	unPackageUTXOS := utxoSet.FindUnPackageSpentableUTXOS(from, txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0
	for _, UTXO := range unPackageUTXOS {
		money += UTXO.Output.Value
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.Index)
		if money >= amount {
			return money, spentableUTXO
		}
	}

	//钱还不够
	utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
		utxobreak:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutPuts(v)
				for _, utxo := range txOutputs.UTXOS {
					money += utxo.Output.Value
					txHash := hex.EncodeToString(utxo.TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)
					if money >= amount {
						break utxobreak
					}
				}
			}
		}
		return nil
	})


	if money < amount {
		log.Panic("余额不足")
	}

	return money, spentableUTXO
}
