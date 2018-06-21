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
		fmt.Printf("当前区块的数据：%s\n", block.PrevBlockHash)
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


//createBlockchain
func (cli *CLI) createBlockchain(address string){
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("DONE!")
}

//getBalance
func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

// send
func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

//Run 参数解析
func (cli *CLI) Run(){
	//检查参数有效性
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address","","The address to send genesis block reward to")

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		// cli.printUsage()
		os.Exit(1)
	}

	//判断命令是否被调用
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == ""{
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if sendCmd.Parsed(){
		if *sendFrom == "" || *sendTo == "" ||*sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	// if printChainCmd.Parsed() {
	// 	cli.printChain()
	// }



}