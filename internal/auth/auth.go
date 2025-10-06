package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("error while creating a password hash - %w\n", err)
	}

	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	isCorrect, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("error while checking if password is valid - %w\n", err)
	}

	return isCorrect, nil
}
