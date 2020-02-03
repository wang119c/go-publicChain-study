package BLC

type TXInput struct {
	//1.交易id
	Txhash []byte
	//2.存储txoutput 在vout 里面的索引
	Vout int
	//3.用户名
	ScriptSig string
}

//判断当前的消费是谁的钱
func (txInput *TXInput) UnlockWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
