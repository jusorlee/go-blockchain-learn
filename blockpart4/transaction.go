package main

import (
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
)

//奖励
const subsidy = 10

// Transaction 表示一个比特币的交易
type Transaction struct {
	ID []byte
	Vin	[]TXInput
	Vout []TXOutput
}

//IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//TXInput 一个交易的输入
type TXInput struct {
	//存储输出所属的交易的ID
	Txid []byte
	//Vout存储输出的序号（ 一个交易可以包括多个TXO）
	Vout int
	//ScriptSig存储一个脚本
	ScriptSig string
}

//TXOutput 一个交易的输出
type TXOutput struct {
	// 存储货币信息
	Value int
	// ScriptPubKey仅仅存储用户定义的字符串
	ScriptPubKey string
}



//SetID 设置一个交易的ID
func (tx *Transaction) SetID(){
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// NewCoinbaseTX 创建一个新的coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == ""{
		data = fmt.Sprintf("奖励给 '%s'", to)
	}
	//输入信息
	// coinbase交易是一种特殊的交易，该TXI不会引用任何TXO，而会直接生成一个TXO，这是作为奖励给矿工的。
	txin := TXInput{[]byte{}, -1, data}
	//输出信息
	// subsidy是挖矿的奖励值，这里在前面设置了全局变量10.
	// 在比特币中，每挖出21000个block是，奖励值减半。
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx

}


// CanUnlockOutputWith 确认地址是否已经初始化交易
func (in *TXInput) CanUnlockOutputWith(unlockData string) bool {
	return in.ScriptSig == unlockData
}


// CanBeUnlockedWith 检查是否可用所提供的数据解锁output
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}


// NewUTXOTransaction 创建一个新的交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	// 返回：累加的交易额acc和UTXO列表validOutputs
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount{
		log.Panic("错误：金额不足")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs{
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	// 找零
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}