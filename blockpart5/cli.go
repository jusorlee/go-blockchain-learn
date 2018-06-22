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
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
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
	if !ValidateAddress(address) {
		log.Panic("Error:Address is not valid")
	}
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("DONE!")
}

//getBalance
func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("Error:Address is not valid")
	}
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	UTXOs := bc.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

// send
func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("Error:Sender Address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("Error:Recipient Address is not valid")
	}

	bc := NewBlockchain(from)
	defer bc.db.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

// createWallet创建钱包
func (cli *CLI) createWallet(){
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()

	wallets.SaveToFile()
	fmt.Printf("你的新地址：%s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}

	address := wallets.GetAddresses()
	fmt.Println(address)

	for _, address := range address {
		fmt.Println(address)
	}
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

	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)

	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
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

	if createWalletCmd.Parsed(){
		cli.createWallet()
	}

	if listAddressesCmd.Parsed(){
		cli.listAddresses()
	}

	// if printChainCmd.Parsed() {
	// 	cli.printChain()
	// }



}