package main

//在这个文件里面，是把每一个区块连起来，怎么连？
//每一个区块，都包含上一个区块的hash,就是使用这个hash将每个区块进行链接。

import (
	"fmt"
	bolt "github.com/coreos/bbolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

//Blockchain 定义一个机构体，来存储这条链
type Blockchain struct {
	tip []byte
	db *bolt.DB
}

// BlockchainIterator 迭代所有区块链上的block
type BlockchainIterator struct{
	currentHASH []byte
	db *bolt.DB
}

//NewBlockchain 创建一个新的区块链
func NewBlockchain() *Blockchain {
	var tip []byte

	//打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// 这里是函数的核心。在这里，我们先获取了存储区块的 bucket：
		//如果存在，就从中读取 l 键；
		//如果不存在，就生成创世块，创建 bucket，并将区块保存到里面，
		//然后更新 l 键以存储链中最后一个块的哈希。
		//Bucket 可以暂时理解为表格
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("区块链不存在，正在创建...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("1"), genesis.Hash)
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

//AddBlock 在区块链中添加一个区块
func (bc *Blockchain) AddBlock(data string) {
	/*
	//前一个区块是当前区块的长度减1
	prevBlock := bc.blocks[len(bc.blocks)-1]
	fmt.Println("prevBlock:", prevBlock)

	//当前区块的设置
	newblock := NewBlock(data, prevBlock.Hash)
	fmt.Println("newblock:", newblock)

	//将新生成的区块，添加到列表中
	bc.blocks = append(bc.blocks, newblock)
	fmt.Println("bc.blocks:", bc.blocks)
	*/
	var lastHash []byte

	//从数据库里面取出最后一个区块的hash
	err := bc.db.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error{
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
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

//Next 返回下一个区块信息
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHASH)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHASH = block.PrevBlockHash
	fmt.Println(i.currentHASH)
	return block
}

//这个main函数也是调试用，现在我们把它抽取到一个单独的文件中
/*
func main() {
	bc := NewBlockchain()
	fmt.Println(bc)
	bc.AddBlock("test1")
}
*/
