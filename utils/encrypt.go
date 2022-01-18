package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	compare := bcrypt.CompareHashAndPassword(existing, incoming)
	return compare == nil
}
