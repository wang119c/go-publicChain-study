package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

const dbName = "blockChain.db"  //数据库名称
const blockTableName = "blocks" //表名

type Blockchain struct {
	//Blocks []*Block // 存储有序的区块

	Tip []byte //最新的区块的hash
	DB  *bolt.DB
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
				fmt.Printf("%x\n", in.Txid)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%s\n", out.ScriptPubKey)
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
func CreateBlockchainWithGenesisBlock(address string) {
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
func MineNewBlock(from []string, to []string, amount []string) {
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)
}
