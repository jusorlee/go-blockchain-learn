package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

//难度系数
//算出来的HASH前24位必须是0，如果用16进制表示，前6位必须是0
//targetBits 的值越大，计算时间越长
const targetBits = 20

//ProofOfWork 是一个proof-of-work的类型
type ProofOfWork struct {
	//一个区块指针
	block *Block
	//目标指针，“target”表示困难度，其类型为big.Int，之所以使用big结构是为了方便和hash值做比较检查是否满足要求。
	target *big.Int
}

//NewProofOfWork 创建一个新的ProofOfWork，并返回，这里只是对target进行处理
func NewProofOfWork(b *Block) *ProofOfWork {
	//设置big.NewInt的初始化位1
	target := big.NewInt(1)
	fmt.Println("target:", target)
	//结果：target: 1
	//左移256-targetBits位，也就是
	target.Lsh(target, uint(256-targetBits))
	fmt.Println("target.Lsh 256:", target)
	//结果：target.Lsh 256: 6901746346790563787434755862277025452451108972170386555162524223799296

	//传进来的block原封不动，和target打包返回
	pow := &ProofOfWork{b, target}
	return pow
}

//prePareData 把一个区块里面的信息链接起来
func (pow *ProofOfWork) prePareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	//	fmt.Println("prePareData->data:", data)
	return data
}

//Run proof of work的核心算法，挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("挖矿：正在计算\"%s\"的hash值\n", pow.block.Data)
	fmt.Println("maxNonce：", maxNonce)
	for nonce < maxNonce {
		data := pow.prePareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		// hash转换成Big Integer,方便做比较
		hashInt.SetBytes(hash[:])

		// 比较hashInt和target的大小。hashInt<target时返回-1；hashInt>target时返回+1；否则返回0。
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}

	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

//Validate 验证block's pow
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prePareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
