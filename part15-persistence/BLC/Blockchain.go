package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

const dbName = "blockChain.db"  //数据库名称
const blockTableName = "blocks" //表名

type Blockchain struct {
	//Blocks []*Block // 存储有序的区块

	Tip []byte //最新的区块的hash
	DB  *bolt.DB
}

//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string, height int64, preHash []byte) () {
	//增加新区块
	newBlock := NewBlock(data, height, preHash)
	//往链中增加区块
	blc.Blocks = append(blc.Blocks, newBlock)
}

//1.创建带创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic("blocks table create failed")
		}

		if b == nil {
			//创建创世区块
			genesiBlock := CreateGenesisBlock("Genesis Data...")
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

	//返回区块链对象
	return &Blockchain{blockHash, db}
}
