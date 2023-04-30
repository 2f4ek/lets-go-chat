package repositories

import (
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
)

var users = make(map[string]models.User)

func AppendUser(user models.User) {
	users[user.Name] = user
}

func CreateUser(userName string, userPassword string) (*models.User, bool) {
	_, userExists := users[userName]
	passwordHash, _ := hasher.HashPassword(userPassword)

	return &models.User{
		Name:     userName,
		Password: passwordHash,
		Token:    helpers.GenerateSecureToken(),
	}, userExists
}

func GetUser(userName string) (models.User, bool) {
	user, userExists := users[userName]
	return user, userExists
}
