package main

import (
	"fmt"
)

//Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// AddBlock 在区块链中保存私有数据
func (bc *Blockchain) AddBlock(data string) {
	//前一个区块是当前区块的长度-1
	prevBlock := bc.blocks[len(bc.blocks)-1]
	fmt.Println("prevBlock:",prevBlock)
	
	//创建新的区块
	newBlock := NewBlock(data,prevBlock.Hash)
	fmt.Println("newBlock:",newBlock)

	//把新建的区块添加到原区块链表的后面
	bc.blocks = append(bc.blocks, newBlock)
	fmt.Println(bc.blocks)
}

// NewBlockchain 创建一个新的区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

