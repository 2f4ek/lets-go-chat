package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

const password = "password"
const wrongPassword = "thisPasswordIsWrong"
const tooLongPassword = "passwordShouldHaveMoreThan72BytesPasswordShouldHaveMoreThan72BytesAndThisIsIt"

func TestHashPassword(t *testing.T) {
	hash, _ := HashPassword(password)
	if result := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); result != nil {
		t.Errorf("Hashed passwords missmatch")
	}
}

func TestHashPasswordFailure(t *testing.T) {
	if _, err := HashPassword(tooLongPassword); err == nil {
		t.Errorf("Max password length should be 72 bytes")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	bcryptHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if status := CheckPasswordHash(password, string(bcryptHash)); status != true {
		t.Errorf("Hashed passwords mismatch")
	}
}

func TestCheckPasswordHashFailure(t *testing.T) {
	hash, _ := HashPassword(password)
	if status := CheckPasswordHash(wrongPassword, hash); status == true {
		t.Errorf("Password hashes should mismatch")
	}
}
