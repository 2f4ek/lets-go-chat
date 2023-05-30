package repositories_test

import (
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUser(t *testing.T) {
	userName := "TestGetUser"
	user, _ := repositories.CreateUser(userName, "password1")

	assert.IsType(t, &models.User{}, user)

	repositories.AppendUser(*user)
	user, _ = repositories.GetUser(userName)

	assert.IsType(t, &models.User{}, user)
}

func TestGetUserByToken(t *testing.T) {
	userName := "TestGetUserByToken"
	token := helpers.GenerateSecureToken()
	user, _ := repositories.CreateUser(userName, "password1")
	user.Token = token

	assert.IsType(t, &models.User{}, user)

	repositories.AppendUser(*user)
	user = repositories.GetUserByToken(token)

	assert.IsType(t, &models.User{}, user)
}

func TestUpdateToken(t *testing.T) {
	userName := "TestUpdateToken"
	token1 := helpers.GenerateSecureToken()
	token2 := helpers.GenerateSecureToken()
	user, _ := repositories.CreateUser(userName, "password1")
	user.Token = token1

	assert.Equal(t, user.Token, token1)

	repositories.AppendUser(*user)
	repositories.UpdateToken(user, token2)

	assert.Equal(t, user.Token, token2)
}

func TestRevokeToken(t *testing.T) {
	userName := "TestRevokeToken"
	token := helpers.GenerateSecureToken()
	user, _ := repositories.CreateUser(userName, "password1")
	user.Token = token
	repositories.AppendUser(*user)
	repositories.RevokeToken(user)

	assert.Empty(t, user.Token)
	assert.Empty(t, repositories.GetUserByToken(token))
}

func TestUserExistsError(t *testing.T) {
	userName := "TestUserExistsError"
	user, _ := repositories.CreateUser(userName, "password1")
	repositories.AppendUser(*user)
	_, err := repositories.CreateUser(userName, "password1")

	assert.NotEmpty(t, err)
}
