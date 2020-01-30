package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		//创建Block表
		b, err := tx.CreateBucket([]byte("Block"))
		if err != nil {
			return fmt.Errorf("create bucket:%s", err)
		}
		//往表里面存储数据
		if b != nil {
			err := b.Put([]byte("l"), []byte("send 100 btc to huizi"))
			if err != nil {
				log.Panic("数据出错")
			}
		}
		//返回nil,以便数据库操作
		return nil
	})

	//更新失败
	if err != nil {
		log.Panic(err)
	}

}
