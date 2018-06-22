package main

//在这个文件里面，是把每一个区块连起来，怎么连？
//每一个区块，都包含上一个区块的hash,就是使用这个hash将每个区块进行链接。

import (
	"fmt"
	bolt "github.com/coreos/bbolt"
	"log"
	"os"
	"encoding/hex"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The time 10/jun/2018 chancellor on brink of second  bailout for banks"

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
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found.Create one first")
		os.Exit(1)
	}

	var tip []byte

	//打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("1"))
		
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	
	bc := Blockchain{tip, db}

	return &bc
}

//dbExists 判断区块链数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile);os.IsNotExist(err){
		return false
	}
	return true
}

// CreateBlockchain 创建一个新的区块链数据库DB
func CreateBlockchain(address string) *Blockchain  {
	if dbExists(){
		fmt.Println("区块链已经存在。")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// 调用NewCoinbaseTx 获取交易信息
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		// 生成创世块
		genesis := NewGenesisBlock(cbtx)

		//创建Bucket
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		//往Bucket里面写入数据
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		//把最上面的一个区块的hash绑定为1
		err = b.Put([]byte("1"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
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


// 查找包含UTXO的交易
// FindUnspentTransactions 返回一个包含未花费输出交易列表
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	//未花费交易
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions{
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout{
				//Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in :=range tx.Vin{
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

//FindUTXO 找出并返回所有未花费交易输出
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

//FindSpendableOutputs 找到并返回未花费输出引用到输入，传入的参数是地址，所需金额
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	// 找到所有未花费的交易
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	// 遍历账户所有的UTX
	for _, tx := range unspentTXs{
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Vout {
			// 找到账户下的输出，并判断交易额是否满足需求
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				// 不满足需求，则累加交易额
				accumulated += out.Value
				// 加入UTXO的列表中
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				// 累加的交易额，大于等于需求，则停止遍历
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	// 返回累加交易额accumulated和UTXO列表unspentOutputs
	return accumulated, unspentOutputs
}



// MineBlock 挖矿
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

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


//这个main函数也是调试用，现在我们把它抽取到一个单独的文件中
/*
func main() {
	bc := NewBlockchain()
	fmt.Println(bc)
	bc.AddBlock("test1")
}
*/
