package chatModels

import (
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/gorilla/websocket"
)

type IChatUser interface {
	GetUserId() models.UserId
	ReadMessage()
	WriteMessage()
	SyncMissedMessages()
}

type ChatUser struct {
	Conn          *websocket.Conn
	Chat          *Chat
	MessageChanel chan []byte
	User          models.User
}

func (cu *ChatUser) GetUserId() models.UserId {
	return cu.User.Id
}

func (cu *ChatUser) ReadMessage() {
	defer func() {
		cu.Chat.Logout <- *cu
		err := cu.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, p, err := cu.Conn.ReadMessage()
		if err != nil {
			break
		}

		message := models.Message{Message: string(p)}
		_, err = cu.Chat.ChatMessageRep.Save(message)
		if err != nil {
			break
		}

		cu.Chat.MessageBroadcast <- p
	}
}

func (cu *ChatUser) WriteMessage() {
	for {
		select {
		case message := <-cu.MessageChanel:
			err := cu.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func (cu *ChatUser) SyncMissedMessages() {
	messages, err := cu.Chat.ChatMessageRep.GetMissedMessages(&cu.User)
	if err != nil {
		return
	}

	for _, message := range messages {
		err := cu.Conn.WriteMessage(websocket.TextMessage, []byte(message.Message))

		if err != nil {
			return
		}
	}
}
