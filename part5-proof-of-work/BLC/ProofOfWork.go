package BLC

type ProofOfWork struct {
	Block *Block //要验证的区块
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {

	



	return nil, 0
}

//创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	return &ProofOfWork{block}
}
