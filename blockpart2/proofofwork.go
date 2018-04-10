package main


import (
	"fmt"
	"math"
	"math/big"
	"bytes"
	"crypto/sha256"
)

var (
	maxNonce = math.MaxInt64
)

// 难度系数
//算出来的哈希前 24 位必须是 0，如果用 16 进制表示，就是前 6 位必须是 0
//targetBits 的值越大，计算时间越长。
const targetBits =24

// ProofOfWork 是一个 proof-of-work的类型
type ProofOfWork struct {
	block *Block	//一个区块
	target *big.Int	//目标指针
}

// NewProofOfWork 创建一个新的ProofOfWork,并返回
func NewProofOfWork(b *Block) *ProofOfWork{
	// 设置big.NewInt的初始化为1
	target := big.NewInt(1)
	fmt.Println("target1:",target)
	// 左移 256 - targetBits 位
	target.Lsh(target, uint(256-targetBits))
	// 为以后的值
	fmt.Println("target2:",target)
	pow := &ProofOfWork{b, target}
	fmt.Println(pow)
	return pow
}

// prepareData 把区块里的信息链接起来
func (pow *ProofOfWork) prepareData(nonce int) []byte{
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
	return data
}
//Run  运行 proof of work
func (pow *ProofOfWork) Run() (int, []byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("挖矿：正在计算 \"%s\" 的Hash值\n",pow.block.Data)
	fmt.Println("maxNonce:",maxNonce)
	for nonce < maxNonce{
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		// hash转换成Big Integer
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

// Validate 验证 block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	// HashInt 是 hash 的整形表示
	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}