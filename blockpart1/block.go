package main

import (
	"fmt"
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
)

// Block 结构体是区块信息
type Block struct {
	Timestamp     int64		//当前时间戳，也就是区块创建的时间
	Data          []byte	//区块存储的实际有效信息，也就是交易
	PrevBlockHash []byte	//前一个块的哈希，即父哈希
	Hash          []byte	//当前块的哈希
}

// SetHash 方法计算和设置区块的hash
func (b *Block) SetHash(){
	//strconv.FormatInt 将一个字符串解析为整数
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	fmt.Println("timestamp:", timestamp)
	
	// 测试bytes.Join()
	// fmt.Println(bytes.Join([][]byte{{'1'},{'3'},{'3'}},[]byte{}))
	
	// headers 头部信息，把前一个区块hash,当前Data,当前时间戳链接起来
	headers := bytes.Join([][]byte{b.PrevBlockHash,b.Data,timestamp}, []byte{})
	fmt.Println("headers:",headers)
	
	// 当前hash计算，把headers的组合信息，加密
	hash := sha256.Sum256(headers)
	fmt.Println("hash:",hash)
	
	//设置当前hash值
	b.Hash = hash[:]
	
}

// NewBlock creates and returns Block
func NewBlock(data string, PrevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), PrevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// NewGenesisBlock creates and returns genesis Block,is the first BLock
func NewGenesisBlock() *Block {
	return NewBlock("Hello,World",[]byte{})
}
