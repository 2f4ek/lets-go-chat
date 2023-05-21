package repositories

import (
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
	uuid "github.com/satori/go.uuid"
)

var users = make(map[string]models.User)

func AppendUser(user models.User) {
	users[user.Name] = user
}

func CreateUser(userName string, userPassword string) (*models.User, bool) {
	_, userExists := users[userName]
	if userExists {
		return nil, userExists
	}

	passwordHash, _ := hasher.HashPassword(userPassword)

	return &models.User{
		Id:       uuid.NewV4().String(),
		Name:     userName,
		Password: passwordHash,
		Token:    helpers.GenerateSecureToken(),
	}, userExists
}

func GetUser(userName string) (models.User, bool) {
	user, userExists := users[userName]
	return user, userExists
}

func GetUserByToken(token string) *models.User {
	for _, user := range users {
		if user.Token == token {
			return &user
		}
	}

	return nil
}

func UpdateToken(user models.User, token string) {
	user.Token = token
	users[user.Name] = user
}

func RevokeToken(user *models.User) {
	actualUser := users[user.Name]
	actualUser.Token = ""
	users[user.Name] = actualUser
}
