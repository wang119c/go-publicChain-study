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
func (blc *Blockchain) AddBlockToBlockchain(data string) () {
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
			newBlock := NewBlock(data, block.Height+1, block.Hash)
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
		b, err := tx.CreateBucketIfNotExists([]byte(blockTableName))
		if err != nil {
			log.Panic("blocks table create failed")
		}

		if b != nil {
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
