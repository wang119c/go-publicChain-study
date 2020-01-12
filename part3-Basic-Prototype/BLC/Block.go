package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块HASH
	PrevBlockHash []byte
	//3.交易数据
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.Hash
	Hash []byte
}

func (block *Block) setHash() {
	//1.height []byte
	heightBytes := IntToHex(block.Height)
	//2.将时间戳转[]byte
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timeBytes := []byte(timeString)
	//3.拼接所有属性
	headers := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes, block.Hash}, []byte{})
	//4.生成hash
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

//创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {

	//创建区块
	block := &Block{
		height,
		prevBlockHash,
		[]byte(data),
		time.Now().Unix(),
		nil,
	}
	//设置Hash
	block.setHash()

	return block
}

//单独写一个方法生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
