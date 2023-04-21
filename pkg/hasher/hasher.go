package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

const cost = 10

func HashPassword(password string) (string, error) {
	var bytes, err = bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal("Server error while password hash")
	}

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
