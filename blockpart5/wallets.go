package main


import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 保存的本地文件
const walletFile = "wallet.dat"

// Wallets 存储手机的钱包
type Wallets struct {
	Wallets map[string]*Wallet
}

// LoadFromFile 从文件中加载wallets
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}


// NewWallets 创建wallets和从一个文件中读取
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// CreateWallet 创建Wallet
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s",wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetAddress 返回一个地址数组
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses,address)
	}

	return addresses
}

// GetWallet 根据地址，返回Wallet
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// SaveToFile 保存wallets到一个文件里面
func (ws Wallets) SaveToFile(){
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)

	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}