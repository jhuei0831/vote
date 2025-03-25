package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher 定義了一個雜湊接口，包含 HashPassword 和 ComparePassword 方法。
type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) bool
}

// SHA256Hasher 是一個實現了 Hasher 接口的結構體，使用 SHA-256 進行雜湊。
type SHA256Hasher struct{}

// HashPassword 使用 bcrypt 演算法將密碼進行雜湊處理。
// 傳入的密碼會被轉換為 byte slice，並使用成本參數 14 進行雜湊。
// 返回雜湊後的密碼字串和可能發生的錯誤。
func (ph *SHA256Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// ComparePassword 比較給定的密碼和雜湊密碼是否匹配。
// 參數:
// - password: 原始密碼的字串。
// - hashedPassword: 雜湊後的密碼字串。
// 回傳值:
// - 如果密碼匹配則回傳 true，否則回傳 false。
func (ph *SHA256Hasher) ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
