package main

import (
	"fmt"
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

//Block结构体是区块链信息
type Block struct {
	//当前时间戳
	Timestamp int64
	//区块的的有消息
	Data []byte
	//前一个区块的HASH
	PrevBlockHash []byte
	//当前块的Hash
	Hash []byte
	//计数器
	Nonce int
}

//NewBlock 创建新的区块，并返回区块信息
func NewBlock(data string, prevBlockHash []byte) *Block {
	fmt.Println("正在新增一个区块")
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}

	//pow计算
	pow := NewProofOfWork(block)
	fmt.Println("pow:", pow)

	//挖矿
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	//打印时间戳
	fmt.Println("Timestamp:", block.Timestamp)
	//打印交易信息
	fmt.Println("Data:", block.Data)
	//打印上一个区块的hash
	fmt.Println("PrevBlockHash:", block.PrevBlockHash)
	//打印当前hash
	fmt.Println("Hash:", block.Hash)
	//返回区块信息
	return block
}

//NewGenesisBlock创世块的生成,也就是第一个区块的生成
func NewGenesisBlock() *Block {
	fmt.Println("正在生成创世块....请耐心等待")
	return NewBlock("这是第一个区块，也叫创世块", []byte{})
}

//Serialize 序列化区块信息，才可以在boltDB中使用
func (b *Block) Serialize() []byte {
	//定义一个buffer，存储序列化之后的数据
	var result bytes.Buffer
	//初始化一个 go encoder
	encoder := gob.NewEncoder(&result)
	//对block进行编码
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeserializeBlock 解序列化
func DeserializeBlock(d []byte) *Block {
	//定义一个Block 存储解序列化之后的数据
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//这个是区块链的第一个文件代码，所以使用main函数来测试，现在测试完成，我们注释掉。
//下面我们开始写第二个文件：把所有每一个区块连起来，变成一个列表
/*
func main() {
	NewGenesisBlock()
	NewBlock("hello", []byte{})

}
*/
