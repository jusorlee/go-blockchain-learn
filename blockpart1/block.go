package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
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
}

//SetHash 计算和设置区块的hash
func (b *Block) SetHash() {
	//strconv.FormatInt  将一个字符串解析为整数
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	fmt.Println(timestamp)
	//Println的结果：[49 53 50 55 53 55 48 57 52 49]

	//测试bytes.Join()
	//fmt.Println(bytes.Join([][]byte{{'a'}, {'3'}, {'3'}}, []byte{}))
	//结果：[97 51 51]

	//header 头部信息，把前一个区块的hash,当前Data,当前时间戳Join在一起
	header := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	fmt.Println("header:", header)
	//结果：header: [104 101 108 108 111 49 53 50 55 53 55 49 53 56 53]

	//当前的hash计算，把header的组合信息加密
	hash := sha256.Sum256((header))
	fmt.Println("hash:", hash)
	//结果：hash: [157 215 223 149 65 120 46 89 107 86 202 88 117 30 155 200 174 121 205 116 100 144 103 151 37 137 135 214 213 100 149 241]

	//设置当前区块的hash值
	b.Hash = hash[:]

}

//NewBlock 创建新的区块，并返回区块信息
func NewBlock(data string, prevBlockHash []byte) *Block {
	fmt.Println("正在新增一个区块")
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	//设置当前区块的hash
	block.SetHash()

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

//这个是区块链的第一个文件代码，所以使用main函数来测试，现在测试完成，我们注释掉。
//下面我们开始写第二个文件：把所有每一个区块连起来，变成一个列表
/*
func main() {
	NewGenesisBlock()
	NewBlock("hello", []byte{})

}
*/
