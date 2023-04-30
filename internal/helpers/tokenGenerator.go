package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

var length = 16

func GenerateSecureToken() string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
