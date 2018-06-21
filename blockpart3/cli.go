package main

import (
	"fmt"
	"os"
	"strconv"
	"flag"
	"log"
)

// CLI 是命令行参数
type CLI struct {
	bc *Blockchain
}

//printUsage 提示
func (cli *CLI) printUsage(){
	fmt.Println("Usage: ")
	fmt.Println(" addblock -data BLOCK_DATA -向区块链中添加一个区块")
	fmt.Println(" printchain -打印所有区块信息")
}

//validateArgs 检查用户输入的参数
func (cli *CLI) validateArgs(){
	if len(os.Args) < 2 {
		fmt.Println("您输入的参数不正确！！！")
		cli.printUsage()
		os.Exit(1)
	}
}

//printChain 打印数据链信息
func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("当前区块的hash:%x\n", block.Hash)
		fmt.Printf("当前区块的数据：%s\n", block.Data)
		fmt.Printf("计算次数：%d\n", block.Nonce)
		fmt.Println()
		pow := NewProofOfWork(block)

		fmt.Printf("pow:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}
}

func (cli *CLI) addBlock (data string){
	cli.bc.AddBlock(data)
	fmt.Println("成功！")
}

//Run 参数解析
func (cli *CLI) Run(){
	//检查参数有效性
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data","","Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//判断命令是否被调用
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}



}