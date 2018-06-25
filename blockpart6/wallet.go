package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet 存储私钥和公钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}


// newKeyPair 返回公钥和私钥
// 函数使用椭圆曲线算法生成私钥,紧接着通过私钥生成公钥
// 需要注意一点：椭圆曲线算法中，公钥是曲线上的点集合
// 公钥由X,Y坐标混合而成
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// NewWallet 创建并返回一个钱包

func NewWallet() *Wallet {
	
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	fmt.Println(wallet)

	return &wallet
}

// HashPubKey hash公钥
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160

}


// Checksum 生成一个checksum的公钥
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)

	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

 
// GetAddress 返回钱包地址
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionPayload := append([]byte{version}, pubKeyHash...)

	checksum := Checksum(versionPayload)

	fullPayload := append(versionPayload, checksum...)

	address := Base58Encode(fullPayload)

	return address

}

// ValidateAddress 验证地址的有效性
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1: len(pubKeyHash) - addressChecksumLen]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// func main(){
// 	NewWallet()
// }