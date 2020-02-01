package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块HASH
	PrevBlockHash []byte
	//3.交易数据
	Txs []*Transaction

	//4.时间戳
	Timestamp int64
	//5.Hash
	Hash []byte
	//6.Nonce  工作证明
	Nonce int64
}

//需要将Txs转换为字节数组
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//序列化字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {

	//创建区块
	block := &Block{
		height,
		prevBlockHash,
		txs,
		time.Now().Unix(),
		nil,
		0,
	}

	//调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	//挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Println(" ")
	return block
}

//单独写一个方法生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
