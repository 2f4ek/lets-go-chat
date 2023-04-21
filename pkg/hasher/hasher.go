package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword function is used to hash a user's password for security purposes.
func HashPassword(password string) (string, error) {
	var bytes, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

// CheckPasswordHash function is used to check if a given password matches a given hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
