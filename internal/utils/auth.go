package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	var err error
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(pass), nil
}

func VerifyPassword(password string, validate string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(validate))

	return err == nil
}
