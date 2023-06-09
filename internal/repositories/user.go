package repositories

import (
	"time"

	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
)

type IUserRepository interface {
	AppendUser(user models.User)
	CreateUser(userName string, userPassword string) (*models.User, bool)
	GetUser(userName string) (*models.User, bool)
	GetUserByToken(token string) *models.User
	UpdateToken(user *models.User, token string)
	RevokeToken(user *models.User)
	UpdateUserLastActivity(user *models.User)
}

type UserRepository struct {
	users map[models.UserId]models.User
}

func ProvideUserRepository() *UserRepository {
	return &UserRepository{users: make(map[models.UserId]models.User)}
}

func (rep *UserRepository) AppendUser(user models.User) {
	rep.users[user.Id] = user
}

func (rep *UserRepository) CreateUser(userName string, userPassword string) (*models.User, bool) {
	for _, user := range rep.users {
		if user.Name == userName {
			return nil, true
		}
	}

	passwordHash, _ := hasher.HashPassword(userPassword)

	userId := len(rep.users)
	userId++

	return &models.User{
		Id:       models.UserId(userId),
		Name:     userName,
		Password: passwordHash,
		Token:    helpers.GenerateSecureToken(),
	}, false
}

func (rep *UserRepository) GetUser(userName string) (*models.User, bool) {
	for _, user := range rep.users {
		if user.Name == userName {
			return &user, true
		}
	}

	return nil, false
}

func (rep *UserRepository) GetUserByToken(token string) *models.User {
	for _, user := range rep.users {
		if user.Token == token {
			return &user
		}
	}

	return nil
}

func (rep *UserRepository) UpdateToken(user *models.User, token string) {
	user.Token = token
	rep.users[user.Id] = *user
}

func (rep *UserRepository) RevokeToken(user *models.User) {
	user.Token = ""
	rep.users[user.Id] = *user
}

func (rep *UserRepository) UpdateUserLastActivity(user *models.User) {
	user.LastActivity = time.Now()
	rep.users[user.Id] = *user
}
