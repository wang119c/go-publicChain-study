package BLC

type Blockchain struct {
	Blocks []*Block // 存储有序的区块
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
	//创建创世区块
	genesiBlock := CreateGenesisBlock("Genesis Data...")
	//返回区块链对象
	return &Blockchain{[]*Block{genesiBlock}}
}
