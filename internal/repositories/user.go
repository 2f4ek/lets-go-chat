package repositories

import (
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
	"time"
)

var users = make(map[models.UserId]models.User)

func AppendUser(user models.User) {
	users[user.Id] = user
}

func CreateUser(userName string, userPassword string) (*models.User, bool) {
	for _, user := range users {
		if user.Name == userName {
			return nil, true
		}
	}

	passwordHash, _ := hasher.HashPassword(userPassword)

	userId := len(users)
	userId++

	return &models.User{
		Id:       models.UserId(userId),
		Name:     userName,
		Password: passwordHash,
		Token:    helpers.GenerateSecureToken(),
	}, false
}

func GetUser(userName string) (*models.User, bool) {
	for _, user := range users {
		if user.Name == userName {
			return &user, true
		}
	}

	return nil, false
}

func GetUserByToken(token string) *models.User {
	for _, user := range users {
		if user.Token == token {
			return &user
		}
	}

	return nil
}

func UpdateToken(user *models.User, token string) {
	user.Token = token
	users[user.Id] = *user
}

func RevokeToken(user *models.User) {
	user.Token = ""
	users[user.Id] = *user
}

func UpdateUserLastActivity(user *models.User) {
	user.LastActivity = time.Now()
	users[user.Id] = *user
}
