package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

const dbName = "blockChain.db"  //数据库名称
const blockTableName = "blocks" //表名

type Blockchain struct {
	//Blocks []*Block // 存储有序的区块

	Tip []byte //最新的区块的hash
	DB  *bolt.DB
}

//转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpentableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//1.获取所有的utxo
	utxos := blockchain.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)

	//遍历
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

//遍历输出所有区块信息
func (blc *Blockchain) Printchain() {
	blockchainInterator := blc.Iterator()
	for {
		block := blockchainInterator.Next()

		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)

		fmt.Println("Txs:")

		for _, tx := range block.Txs {
			//fmt.Printf("Txs:%v\n", block.Txs)
			//fmt.Println(tx)

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.Txhash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%x\n", in.Pubkey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%x\n", out.Ripemd160Hash)
			}

		}

		fmt.Println("---------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//1.创建表
		b := tx.Bucket([]byte(blockTableName))
		//2.创建新的区块
		if b != nil {
			//获取上个节点存储的 ， 现获取最新的区块
			blockBytes := b.Get(blc.Tip)
			//进行反序列化
			block := DeserializeBlock(blockBytes)

			//3.增加新区块 , 将区块序列化并且存储到数据库中
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4.更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//5.更新blockchain的tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}
}

//判断创世区块是不是存在
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//1.创建带创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	if dbExists() {
		fmt.Println("创世区块已经存在...")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块...")

	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(blockTableName))
		if err != nil {
			log.Panic("blocks table create failed")
		}

		if b != nil {
			//创建创世区块
			//创建一个 Coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address)
			genesiBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			//将创世区块存储到表中
			err := b.Put(genesiBlock.Hash, genesiBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块hash
			err = b.Put([]byte("l"), genesiBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			blockHash = genesiBlock.Hash

		}

		return nil
	})

	return &Blockchain{blockHash, db}
}

//返回blockchain对象
func BlockchainObject() *Blockchain {
	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	//defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//读取最新区块的hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

//挖掘新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	//1.通过相关算法建立Transaction数组
	var txs []*Transaction

	for index, address := range from {
		//1.建立一笔交易
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	//奖励
	tx := NewCoinbaseTransaction(from[0])
	txs = append(txs, tx)

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	//建立新区块之前进行签名认证

	_txs := []*Transaction{}

	for _, tx := range txs {
		if blockchain.VerifyTransaction(tx, _txs) == false {
			log.Panic("签名失败...")
		}
		_txs = append(_txs, tx)
	}

	//2.建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)
	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("l"), block.Hash)
			blockchain.Tip = block.Hash
		}
		return nil
	})

}

//验证数字签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vins {
		prevTX, err := bc.FindTransaction(vin.Txhash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}
	return tx.Verify(prevTXs)
}

//如果一个地址对应的TXOutput 未花费 ， 那么这个Transaction就应该添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(address))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]

				if in.UnlockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.Txhash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	work1:
		for index, out := range tx.Vouts {
			if out.UnlockScriptPubKeyWithAddress(address) {
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
			if out.UnlockScriptPubKeyWithAddress(address) {

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

	blockIterator := blockchain.Iterator()

	for {
		block := blockIterator.Next()
		fmt.Println(block)
		fmt.Println()

		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			//txhash

			//Vins
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					publicKeyHash := Base58Decode([]byte(address))
					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]

					//是否能够解锁
					if in.UnlockRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.Txhash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}

			//Vouts
		work:
			for index, out := range tx.Vouts {
				if out.UnlockScriptPubKeyWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
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

		//满足条件退出
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOs
}

//查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	untxos := blockchain.UnUTXOs(address, []*Transaction{})
	var amount int64
	for _, utxo := range untxos {
		amount = amount + utxo.Output.Value
	}
	return amount
}

//数字签名
func (blockchain *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, txs []*Transaction) {
	if tx.IsCoinbaseTransaction() {
		return
	}
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vins {
		prevTX, err := blockchain.FindTransaction(vin.Txhash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) FindTransaction(ID []byte, txs []*Transaction) (Transaction, error) {

	for _, tx := range txs {
		if bytes.Compare(tx.TxHash, ID) == 0 {
			return *tx, nil
		}
	}

	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

		//if len(block.PrevBlockHash) == 0 {
		//	break
		//}
	}
	return Transaction{}, nil

}

func (blc *Blockchain) FindUTXOMap() map[string]*TXOutputs {
	blcIterator := blc.Iterator()
	//存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)

	utxoMaps := make(map[string]*TXOutputs)

	for {
		block := blcIterator.Next()

		for i := len(block.Txs) - 1; i >= 0; i-- {

			txOutputs := &TXOutputs{[]*UTXO{}}
			tx := block.Txs[i]
			//txHash := hex.EncodeToString(tx.TxHash)



			if tx.IsCoinbaseTransaction() == false {
				for _, txInput := range tx.Vins {
					txHash := hex.EncodeToString(txInput.Txhash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash], txInput)
				}
			}

			txHash := hex.EncodeToString(tx.TxHash)

		workOutLoop:
			for index, out := range tx.Vouts {
				txInputs := spentableUTXOsMap[txHash]
				if len(txInputs) > 0 {

					isSpent := false

					for _, in := range txInputs {
						if index == in.Vout {
							outPublicKey := out.Ripemd160Hash
							inPublicKey := in.Pubkey

							if bytes.Compare(outPublicKey, Ripemd160Hash(inPublicKey)) == 0 {
								if index == in.Vout {
									isSpent = true
									continue workOutLoop
								}
							}
						}
					}

					if isSpent == false {
						utxo := &UTXO{tx.TxHash, index, out}
						txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
					}

				} else {
					utxo := &UTXO{tx.TxHash, index, out}
					txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
				}
			}
			//设置键值对
			utxoMaps[txHash] = txOutputs
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}
	return utxoMaps
}
