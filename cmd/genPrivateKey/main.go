package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Generate new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	// Get the private key's hexadecimal string (save it securely!)
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	fmt.Println("Private Key:", privateKeyHex)

	// Derive public address from private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("Backend Verifier Address:", publicAddress.Hex())

	// Generate signature
	signature, _ := signMessage("hello", privateKey)
	fmt.Println("Signature:", hex.EncodeToString(signature))

	// Verify signature
	recoveredPubKey, _ := crypto.SigToPub(
		crypto.Keccak256Hash([]byte("hello")).Bytes(),
		signature,
	)

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)
	fmt.Println("Recovered Address:", recoveredAddr.Hex()) // Should equal backendVerifier
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// Calculate Keccak256 hash of the message
	messageHash := crypto.Keccak256Hash([]byte(message))

	// Sign (with Ethereum prefix "\x19Ethereum Signed Message:\n32")
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	// Need to set signature recovery ID to 27 or 28 (EIP-155 compatible)
	// signature[64] += 27
	return signature, nil
}
