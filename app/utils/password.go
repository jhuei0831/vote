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

type Password struct {
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

// GeneratePassword 生成指定長度的隨機密碼。
// 參數:
// - number: 指定生成的密碼數量。
// - length: 指定的密碼長度。
// - format: 指定的密碼格式。
// 返回值:
// - 如果成功生成密碼則返回密碼字串切片，否則返回錯誤。
func (p *Password) GeneratePassword(number int, length int, format string) ([]string, error) {
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
    for i := 0; i < number; i++ {
        password, err := generateRandomString(chars, length)
        if err != nil {
            return nil, err
        }
        passwords[i] = password
    }

    return passwords, nil
}

// generateRandomString 生成指定長度的隨機字串。
// 參數:
// - chars: 指定的字符集。
// - length: 指定的字串長度。
// 返回值:
// - 如果成功生成字串則返回字串，否則返回錯誤。
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

// Encrypt 加密字串
func (p *Password) Encrypt(text string) (string, error) {
	password := os.Getenv("VOTER_PASSWORD_ENCRYPT_KEY")
	// 使用提供的密碼創建一個新的 AES 密碼區塊
	block, err := aes.NewCipher([]byte(password))
	if err != nil {
		return "", err
	}

	// 從 base64 解碼初始化向量 (IV)
	ivDecoded, err := base64.StdEncoding.DecodeString(os.Getenv("VOTER_PASSWORD_ENCRYPT_IV"))
	if err != nil {
		return "", err
	}

	// 創建一個字節切片來保存密文，並在開頭留出 IV 的空間
	ciphertext := make([]byte, aes.BlockSize+len(text))

	// 使用 AES 區塊和 IV 創建一個新的 CTR 流密碼
	stream := cipher.NewCTR(block, ivDecoded)

	// 加密明文並將其存儲在密文切片中，從 IV 之後開始
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// 將密文編碼為 base64 並以字串形式返回
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字串
func (p *Password) Decrypt(text string) (string, error) {
	// 使用提供的密碼創建一個新的 AES 密碼區塊
	block, err := aes.NewCipher([]byte(os.Getenv("VOTER_PASSWORD_ENCRYPT_KEY")))
	if err != nil {
		return "", err
	}

	// 從 base64 解碼初始化向量 (IV)
	ivDecoded, err := base64.StdEncoding.DecodeString(os.Getenv("VOTER_PASSWORD_ENCRYPT_IV"))
	if err != nil {
		return "", err
	}

	// 從 base64 解碼密文
	ciphertext, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	// 檢查密文長度是否小於 AES 區塊大小
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// 使用 AES 區塊和 IV 創建一個新的 CTR 流密碼
	stream := cipher.NewCTR(block, ivDecoded)

	// 解密密文並將結果存儲在相同的密文切片中，從 IV 之後開始
	stream.XORKeyStream(ciphertext[aes.BlockSize:], ciphertext[aes.BlockSize:])

	// 返回解密後的字串，從 IV 之後開始
	return string(ciphertext[aes.BlockSize:]), nil
}