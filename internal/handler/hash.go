package handler

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengembalikan password yang sudah di-hash
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
