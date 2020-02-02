package BLC

type TXInput struct {
	//1.交易id
	Txid []byte
	//2.存储txoutput 在vout 里面的索引
	Vout int
	//3.用户名
	ScriptSig string
}
