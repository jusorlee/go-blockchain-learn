package main

import (
	"fmt"
)

func main(){
	// 创建创世块
	bc := NewBlockchain()
	fmt.Println(bc)

	//创建一个区块
	bc.AddBlock("发送一个以太币给朋友Jusor")
	//创建第二个区块
	bc.AddBlock("发送一个以太币给朋友Leilei")
	
	// for循环，打印出每一个区块的信息。
	for _, block := range bc.blocks{
		fmt.Printf("前一个区块的HASH:%x\n", block.PrevBlockHash)
		fmt.Printf("数据:%s\n",block.Data)
		fmt.Printf("当前区块的HASH:%x\n", block.Hash)
		fmt.Println()
	}


}