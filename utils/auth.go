package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)



func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error hashing password : ", err)
		return "", err
	}

	return string(hash), nil
}

func CheckHashedPassword(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(plainPassword), []byte(hashedPassword))
	return err == nil
}