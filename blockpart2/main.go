package main

/**
单个的区块已经有了，
链条也有了，
现在我们来组装区块链
**/

import (
	"fmt"
	"strconv"
)

func main() {
	//新建一条区块链，并且创建创世块
	bc := NewBlockchain()
	fmt.Println(bc)

	//添加第1个区块
	bc.AddBlock("发送一个信息给朋友")
	//添加第2个区块
	bc.AddBlock("发送一个信息给LEILEI")

	//现在我们来看一下这条区块链有多少个区块，也就是区块链的高度
	//高度肯定是3，区块链也是从0开始计算的
	fmt.Println("高度：", len(bc.blocks))

	//使用for循环，打印每一个区块的信息
	for _, block := range bc.blocks {
		fmt.Printf("当前区块的hash:%x\n", block.Hash)
		fmt.Printf("当前区块的数据：%s\n", block.Data)
		fmt.Printf("计算次数：%d\n", block.Nonce)
		fmt.Println()
		pow := NewProofOfWork(block)

		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
