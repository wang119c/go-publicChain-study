package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	//TxHash    []byte
	//TXOutputs []*TXOutput
	UTXOS []*UTXO

}



//序列化字节数组
func (txOutPuts *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(txOutPuts)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}




//反序列化
func DeserializeTXOutPuts (txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewBuffer(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
