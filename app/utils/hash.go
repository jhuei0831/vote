package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher defines a hash interface that includes HashPassword and ComparePassword methods.
type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) bool
}

// SHA256Hasher is a struct that implements the Hasher interface, using SHA-256 for hashing.
type SHA256Hasher struct{}

// HashPassword hashes the password using the bcrypt algorithm.
// The input password is converted to a byte slice and hashed with cost parameter 14.
// Returns the hashed password string and any error that may occur.
func (ph *SHA256Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword compares whether the given password matches the hashed password.
// Parameters:
// - password: the original password string.
// - hashedPassword: the hashed password string.
// Return value:
// - returns true if the password matches, otherwise returns false.
func (ph *SHA256Hasher) ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
