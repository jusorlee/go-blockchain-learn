package main

import (
	"fmt"
	"log"

	bolt "github.com/coreos/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

//Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	db *bolt.DB
}

// BlockchainIterator 迭代所有区块链上的blocks
type BlockchainIterator struct{
	currentHash []byte
	db *bolt.DB
}

// AddBlock 在区块链中保存私有数据
func (bc *Blockchain) AddBlock(data string) {
	/*
	//前一个区块是当前区块的长度-1
	prevBlock := bc.blocks[len(bc.blocks)-1]
	fmt.Println("prevBlock:",prevBlock)
	
	//创建新的区块
	newBlock := NewBlock(data,prevBlock.Hash)
	fmt.Println("newBlock:",newBlock)

	//把新建的区块添加到原区块链表的后面
	bc.blocks = append(bc.blocks, newBlock)
	fmt.Println(bc.blocks)*/

	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil

	})

}

//Iterator 迭代
func (bc *Blockchain) Iterator() *BlockchainIterator{
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

//Next 返回下一个区块信息
func (i *BlockchainIterator) Next() *Block{
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	fmt.Println(i.currentHash)


	return block
}

// NewBlockchain 创建一个新的区块链
func NewBlockchain() *Blockchain {
	// return &Blockchain{[]*Block{NewGenesisBlock()}}

	var tip []byte
	// 打开数据库，打开一个 BoltDB 文件的标准做法。注意，即使不存在这样的文件，它也不会返回错误。
	db, err := bolt.Open(dbFile,0600,nil)
	if err != nil {
		log.Panic(err)
	}

	// 有两种类型的事务：只读（read-only）和读写（read-write）。
	//这里，打开的是一个读写事务（db.Update(...)），因为我们可能会向数据库中添加创世块。
	err = db.Update(func(tx *bolt.Tx) error{
		// 这里是函数的核心。在这里，我们先获取了存储区块的 bucket：
		//如果存在，就从中读取 l 键；
		//如果不存在，就生成创世块，创建 bucket，并将区块保存到里面，
		//然后更新 l 键以存储链中最后一个块的哈希。
		//Bucket 可以暂时理解为表格
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("区块链不存在，正在创建...")
			genesis := NewGenesisBlock()

			//create a bucket 
			//All keys in a bucket must be unique.
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			//save a key/value pair to a bucket,
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("1"),genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil

	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

