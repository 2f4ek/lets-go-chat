package models_test

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/internal/chatModels"
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

	chatUser := &chatModels.ChatUser{
		Conn:          nil,
		Chat:          chat,
		MessageChanel: make(chan []byte),
		User:          *user,
	}

	err := chat.AddUserToChat(&chatUser.User, connections["user_token"])
	if err != nil {
		return
	}

	assert.Len(t, chat.GetActiveUsers(), 1)

	chat.RemoveUser(chatUser)
}

func TestChat_RemoveUser(t *testing.T) {
	connections := map[string]*websocket.Conn{}
	connections["user_token"] = nil
	user, _ := userRepo.CreateUser("userName", "password1")

	chat := handlers.InitChat()
	chatUser := &chatModels.ChatUser{
		Conn:          nil,
		Chat:          chat,
		MessageChanel: make(chan []byte),
		User:          *user,
	}
	err := chat.AddUserToChat(&chatUser.User, connections["user_token"])
	if err != nil {
		return
	}

	assert.Len(t, chat.GetActiveUsers(), 1)

	chat.RemoveUser(chatUser)

	assert.Len(t, chat.GetActiveUsers(), 0)
}

func TestChat_GetActiveUsers(t *testing.T) {
	connections := map[string]*websocket.Conn{}
	connections["first_user_token"] = nil
	connections["second_user_token"] = nil

	chat := handlers.InitChat()

	user1, _ := userRepo.CreateUser("userName", "password1")
	chatUser1 := &chatModels.ChatUser{
		Conn:          nil,
		Chat:          chat,
		MessageChanel: make(chan []byte),
		User:          *user1,
	}
	err := chat.AddUserToChat(user1, connections["user_token"])
	if err != nil {
		return
	}
	user2, _ := userRepo.CreateUser("userName", "password1")
	chatUser2 := &chatModels.ChatUser{
		Conn:          nil,
		Chat:          chat,
		MessageChanel: make(chan []byte),
		User:          *user2,
	}
	err = chat.AddUserToChat(user2, connections["user_token"])
	if err != nil {
		return
	}

	assert.Len(t, chat.GetActiveUsers(), 2)

	chat.RemoveUser(chatUser1)
	chat.RemoveUser(chatUser2)
}

func BenchmarkChat_CreateAndAddUserToChat(b *testing.B) {
	var rep = &chatModels.Chat{
		ChatUsers: map[models.UserId]*chatModels.ChatUser{},
	}
	for i := 1; i < 10000; i++ {
		user := models.User{
			Id:       models.UserId(i),
			Name:     fmt.Sprint(i, "name"),
			Password: fmt.Sprint(i, "password"),
			Token:    fmt.Sprint(i, "token"),
		}
		err := rep.AddUserToChat(&user, nil)
		if err != nil {
			return
		}
	}

	for n := 0; n < b.N; n++ {
		rep.GetActiveUsers()
	}
}
