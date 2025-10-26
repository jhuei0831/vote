package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// 生成新的私鑰
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// 獲取私鑰的十六進制字符串（妥善保存！）
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	fmt.Println("Private Key:", privateKeyHex)

	// 從私鑰推導出公鑰地址
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("Backend Verifier Address:", publicAddress.Hex())

	// 生成簽名
	signature, _ := signMessage("hello", privateKey)
	fmt.Println("Signature:", hex.EncodeToString(signature))

	// 驗證簽名
	recoveredPubKey, _ := crypto.SigToPub(
			crypto.Keccak256Hash([]byte("hello")).Bytes(),
			signature,
	)

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)
	fmt.Println("Recovered Address:", recoveredAddr.Hex()) // 應等於 backendVerifier
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// 計算消息的 Keccak256 哈希
	messageHash := crypto.Keccak256Hash([]byte(message))
	
	// 簽名（加上 Ethereum 前綴 "\x19Ethereum Signed Message:\n32"）
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	// 需要將簽名的 recovery ID 設置為 27 或 28（EIP-155 兼容）
	// signature[64] += 27
	return signature, nil
}