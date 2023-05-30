package models_test

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	userRepo "github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChat_AddUserToChat(t *testing.T) {
	connections := map[string]*websocket.Conn{}
	connections["user_token"] = nil
	user, _ := userRepo.CreateUser("userName", "password1")

	chat := handlers.InitChat()
	chat.AddUserToChat(user, connections["user_token"])

	assert.Len(t, chat.GetActiveUsers(), 1)

	chat.RemoveUser(user.Token)
}

func TestChat_RemoveUser(t *testing.T) {
	connections := map[string]*websocket.Conn{}
	connections["user_token"] = nil
	user, _ := userRepo.CreateUser("userName", "password1")

	chat := handlers.InitChat()
	chat.AddUserToChat(user, connections["user_token"])

	assert.Len(t, chat.GetActiveUsers(), 1)

	chat.RemoveUser(user.Token)

	assert.Len(t, chat.GetActiveUsers(), 0)
}

func TestChat_GetActiveUsers(t *testing.T) {
	connections := map[string]*websocket.Conn{}
	connections["first_user_token"] = nil
	connections["second_user_token"] = nil

	chat := handlers.InitChat()

	user1, _ := userRepo.CreateUser("userName", "password1")
	chat.AddUserToChat(user1, connections["user_token"])
	user2, _ := userRepo.CreateUser("userName", "password1")
	chat.AddUserToChat(user2, connections["user_token"])

	assert.Len(t, chat.GetActiveUsers(), 2)

	chat.RemoveUser(user1.Token)
	chat.RemoveUser(user2.Token)
}

func BenchmarkChat_CreateAndAddUserToChat(b *testing.B) {
	rep := &models.Chat{
		ChatUsers: map[string]models.ChatUser{},
	}
	for i := 1; i < 10000; i++ {
		user := models.User{
			Id:       fmt.Sprint(i, "id"),
			Name:     fmt.Sprint(i, "name"),
			Password: fmt.Sprint(i, "password"),
			Token:    fmt.Sprint(i, "token"),
		}
		rep.AddUserToChat(&user, nil)
	}

	for n := 0; n < b.N; n++ {
		rep.GetActiveUsers()
	}
}
