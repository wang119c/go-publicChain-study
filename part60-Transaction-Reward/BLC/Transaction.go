package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
)

//交易数据 UTXO
type Transaction struct {
	TxHash []byte      //交易hash
	Vins   []*TXInput  //输入
	Vouts  []*TXOutput //输出
}

//判断当前的交易是否是coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].Txhash) == 0 && tx.Vins[0].Vout == -1
}

//1.Transaction 创建分两种情况
//第一种 创世区块创建时的 Transaction

func NewCoinbaseTransaction(address string) *Transaction {
	//代表消费
	txInput := &TXInput{[]byte{}, -1, nil, []byte("Genesis Data")}
	txOutput := NewTXOutput(10, address)
	//txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//第二种 转账时产生的 Transaction
func NewSimpleTransaction(from string, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {

	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from]

	money, spendableUTXODic := blockchain.FindSpentableUTXOS(from, amount, txs)
	var txInputs []*TXInput
	var txOutputs []*TXOutput
	for txHash, indexArray := range spendableUTXODic {
		for _, index := range indexArray {
			//代表消费
			txHashBytes, _ := hex.DecodeString(txHash)
			txInput := &TXInput{txHashBytes, index, nil, wallet.PublicKey}
			//消费
			txInputs = append(txInputs, txInput)
		}
	}

	//转账
	txOutput := NewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)
	//找零
	txOutput = NewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置hash值
	tx.HashTransaction()

	//进行数字签名
	blockchain.SignTransaction(tx, wallet.PrivateKey)

	return tx
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbaseTransaction() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txhash)].TxHash == nil {
			log.Panic("ERROR:Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTX := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].Pubkey = prevTX.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].Pubkey = nil

		//签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[inID].Signature = signature
	}

}

//拷贝一份新的Transaction 用于签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.Txhash, vin.Vout, vin.Signature, vin.Pubkey})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

//验证交易输入的签名
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}
	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.Txhash)].TxHash == nil {
			log.Panic("ERROR:Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vins {
		prevTX := prevTXs[hex.EncodeToString(vin.Txhash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].Pubkey = prevTX.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].Pubkey = nil

		//私钥ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Pubkey)
		x.SetBytes(vin.Pubkey[:(keyLen / 2)])
		y.SetBytes(vin.Pubkey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}
	return true
}
