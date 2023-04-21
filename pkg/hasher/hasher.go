package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

const cost = 14

func HashPassword(password string) string {
	var bytes, err = bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		log.Fatal("Server error while password hash:")
	}

	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
