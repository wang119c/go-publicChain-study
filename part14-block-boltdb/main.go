package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"go-publicChain-study/part14-block-boltdb/BLC"
	"log"
)

func main() {
	////创世区块
	//blockchain := BLC.CreateBlockchainWithGenesisBlock()
	////新区块
	//blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 200RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//blockchain.AddBlockToBlockchain("Send 300RMB To zhangqiang", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	//
	//
	//fmt.Println(blockchain.Blocks)

	block := BLC.NewBlock("Test", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	//fmt.Printf("%d\n", block.Nonce)
	//fmt.Printf("%x\n", block.Hash)

	//proofOfWork := BLC.NewProofOfWork(block)
	//fmt.Printf("%v", proofOfWork.IsValid())

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		//取block表
		b := tx.Bucket([]byte("blocks"))
		//往表里面存储数据
		if b == nil {
			b, err = tx.CreateBucket([]byte("blocks"))
			if err != nil {
				log.Panic("blocks table create failed")
			}
		}

		err := b.Put([]byte("l"), block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		//返回nil,以便数据库操作
		return nil
	})

	//更新失败
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b != nil {
			blockData := b.Get([]byte("l"))
			//fmt.Printf("%s\n", blockData)

			block = BLC.DeserializeBlock(blockData)
			fmt.Printf("%v\n", block)

		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}
