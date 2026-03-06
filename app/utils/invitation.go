package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"os"
)

type Invitation struct {
}

const (
	TYPE_INT       = "int"
	TYPE_EN        = "en"
	TYPE_MIX       = "mix"
	TYPE_MIX_EXCL  = "mixExcl"
	TYPE_MIX_LOWER = "mixLower"
	TYPE_MIX_UPPER = "mixUpper"
)

var (
	regex09 = "0123456789"
	regexAZ = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	regexaz = "abcdefghijklmnopqrstuvwxyz"
)

// GenerateInvitation generates random passwords of specified length.
// Parameters:
// - number: specifies the number of passwords to generate.
// - length: specifies the password length.
// - format: specifies the password format.
// Returns:
// - Returns a slice of password strings if successful, otherwise returns an error.
func (p *Invitation) GenerateInvitation(number uint, length uint, format string) ([]string, error) {
	if length < 6 {
		length = 6
	}

	var chars []rune
	switch format {
	case TYPE_INT:
		chars = []rune(regex09)
	case TYPE_EN:
		chars = append([]rune(regexaz), []rune(regexAZ)...)
	case TYPE_MIX:
		chars = append([]rune(regex09), append([]rune(regexaz), []rune(regexAZ)...)...)
	case TYPE_MIX_LOWER:
		chars = []rune("23456789abcdefghijkmnpqrstuvwxyz")
	case TYPE_MIX_UPPER:
		chars = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZ")
	case TYPE_MIX_EXCL:
		chars = []rune("23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	default:
		return nil, errors.New("unsupported format")
	}

	passwords := make([]string, number)
	for i := uint(0); i < number; i++ {
		password, err := generateRandomString(chars, int(length))
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}

// generateRandomString generates a random string of specified length.
// Parameters:
// - chars: specifies the character set.
// - length: specifies the string length.
// Returns:
// - Returns the string if successful, otherwise returns an error.
func generateRandomString(chars []rune, length int) (string, error) {
	result := make([]rune, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}

// Encrypt encrypts a string
func (p *Invitation) Encrypt(text string) (string, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher([]byte(os.Getenv("APP_ENCRYPT_KEY")))
	if err != nil {
		return "", err
	}

	// Use a fixed initialization vector (IV)
	iv := []byte(os.Getenv("APP_ENCRYPT_IV"))

	// Create a byte slice to hold the ciphertext and place the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(text))
	copy(ciphertext[:aes.BlockSize], iv)

	// Create a new CTR stream cipher using the AES block and IV
	stream := cipher.NewCTR(block, iv)

	// Encrypt the plaintext and store it in the ciphertext slice, starting after the IV
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// Encode the ciphertext to base64 and return it as a string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string
func (p *Invitation) Decrypt(text string) (string, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher([]byte(os.Getenv("APP_ENCRYPT_KEY")))
	if err != nil {
		return "", err
	}

	// Decode the ciphertext from base64
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	// Check if the ciphertext length is less than the AES block size
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract the initialization vector (IV) from the ciphertext
	iv := ciphertext[:aes.BlockSize]

	// Create a new CTR stream cipher using the AES block and IV
	stream := cipher.NewCTR(block, iv)

	// Decrypt the ciphertext and store the result in the same ciphertext slice, starting after the IV
	stream.XORKeyStream(ciphertext[aes.BlockSize:], ciphertext[aes.BlockSize:])

	// Return the decrypted string, starting after the IV
	return string(ciphertext[aes.BlockSize:]), nil
}
