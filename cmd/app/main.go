package main

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
)

func main() {
	password := "password"
	hashedPassword := hasher.HashPassword(password)

	fmt.Println("Password:", password)
	fmt.Println("Hash:", hashedPassword)

	isMatch := hasher.CheckPasswordHash(password, hashedPassword)
	fmt.Println("Match:   ", isMatch)
}
