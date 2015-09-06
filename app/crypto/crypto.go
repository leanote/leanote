// Package crypto contains two cryptographic functions for both storing and comparing passwords.
package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateHash generates bcrypt hash from plaintext password
func GenerateHash(password string) ([]byte, error) {
	hex := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hex, 10)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

// CompareHash compares bcrypt password with a plaintext one. Returns true if passwords match
// and false if they do not.
func CompareHash(digest []byte, password string) bool {
	hex := []byte(password)
	if err := bcrypt.CompareHashAndPassword(digest, hex); err == nil {
		return true
	}
	return false
}
