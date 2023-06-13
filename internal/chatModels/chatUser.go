package chatModels

import (
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/gorilla/websocket"
)

type ChatUser struct {
	Conn          *websocket.Conn
	Chat          *Chat
	MessageChanel chan []byte
	User          models.User
}

func (cu ChatUser) GetUserId() models.UserId {
	return cu.User.Id
}
