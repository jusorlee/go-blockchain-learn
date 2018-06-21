package main

//在这个文件里面，是把每一个区块连起来，怎么连？
//每一个区块，都包含上一个区块的hash,就是使用这个hash将每个区块进行链接。

import (
	"fmt"
)

//Blockchain 定义一个机构体，来存储这条链
type Blockchain struct {
	//这里就是每个区块组成的切片（slice），这里为什么使用切片而不用数组，
	//因为数组定义好后，是不允许再更改长度，而slice,没有这个要求
	blocks []*Block
}

//NewBlockchain 创建一个新的区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//AddBlock 在区块链中添加一个区块
func (bc *Blockchain) AddBlock(data string) {
	//前一个区块是当前区块的长度减1
	prevBlock := bc.blocks[len(bc.blocks)-1]
	fmt.Println("prevBlock:", prevBlock)

	//当前区块的设置
	newblock := NewBlock(data, prevBlock.Hash)
	fmt.Println("newblock:", newblock)

	//将新生成的区块，添加到列表中
	bc.blocks = append(bc.blocks, newblock)
	fmt.Println("bc.blocks:", bc.blocks)
}

//这个main函数也是调试用，现在我们把它抽取到一个单独的文件中
/*
func main() {
	bc := NewBlockchain()
	fmt.Println(bc)
	bc.AddBlock("test1")
}
*/
