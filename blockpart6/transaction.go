package main

import (
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
	
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
	// 签名
	Signature []byte
	// 公钥
	PubKey []byte
}

//TXOutput 一个交易的输出
type TXOutput struct {
	// 存储货币信息
	Value int
	PubKeyHash []byte
}


// UsesKey 检查地址是否初始化交易
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}


// IsLockedWithKey 检查如果output能被公钥拥有者使用
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput 创建一个新的TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
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
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	//输出信息
	// subsidy是挖矿的奖励值，这里在前面设置了全局变量10.
	// 在比特币中，每挖出21000个block是，奖励值减半。
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx

}


// // CanUnlockOutputWith 确认地址是否已经初始化交易
// func (in *TXInput) CanUnlockOutputWith(unlockData string) bool {
// 	return in.ScriptSig == unlockData
// }




// // CanBeUnlockedWith 检查是否可用所提供的数据解锁output
// func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
// 	return out.ScriptPubKey == unlockingData
// }


// NewUTXOTransaction 创建一个新的交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)

	// 返回：累加的交易额acc和UTXO列表validOutputs
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount, to))
	// 找零
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc - amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	// tx.SetID()
	bc.SignTransaction(&tx,wallet.PrivateKey)

	return &tx
}

// Sign 签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction){
	// 如果是Coinbase交易，则直接返回，不做签名处理
	if tx.IsCoinbase(){
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("错误：前一个交易错误")
		}

	}

	txCopy := tx.TrimmedCopy()

	// 对TrimmedCopy交易中所有的TXI进行遍历
	for inID, vin := range txCopy.Vin{
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}

}

// TrimmedCopy 创建一个交易摘要副本来作签名
func (tx *Transaction) TrimmedCopy() Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout,nil ,nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value,vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Hash 返回Transaction的hash值
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Serialize 返回一个序列Transaction
func (tx Transaction) Serialize() []byte{
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}


// Verify 验证交易输入的签名
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase(){
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR:前一个交易错误")
		}
	}

	txCopy := tx.TrimmedCopy()
	// 创建椭圆曲线用于生成键值对
	curve := elliptic.P256()

	// 对于每个TXI的签名进行验证
	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		// 这个过程和Sign方法是一致的，因为验证的数据需要和签名的数据是一致的
		r := big.Int{}
		s := big.Int{}

		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		// 将椭圆曲线的X,Y坐标点集合（ 其实也是两个字节序列） 组合生成TXI的PubKey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}