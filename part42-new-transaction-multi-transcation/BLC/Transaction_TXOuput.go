package BLC

type TXOutput struct {
	Value        int64  //消费值
	ScriptPubKey string //用户名
}

//判断当前的消费是谁的钱
func (txOutput *TXOutput) UnlockScriptPubKeyWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
