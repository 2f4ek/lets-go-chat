package chatModels

import (
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gorilla/websocket"
)

type Chat struct {
	ChatUsers        map[models.UserId]*ChatUser
	Leavers          chan ChatUser
	MessageBroadcast chan []byte
}

func (c *Chat) RunChat() {
	for {
		select {
		case user := <-c.Leavers:
			if _, ok := c.ChatUsers[user.GetUserId()]; ok {
				c.RemoveUser(&user)
			}
		case message := <-c.MessageBroadcast:
			for _, user := range c.ChatUsers {
				select {
				case user.MessageChanel <- message:
				default:
					c.RemoveUser(user)
				}
			}
		}
	}
}

func (c *Chat) AddUserToChat(user *models.User, conn *websocket.Conn) error {
	u := &ChatUser{
		Conn:          conn,
		Chat:          c,
		MessageChanel: make(chan []byte),
		User:          *user,
	}
	if activeUser, ok := c.ChatUsers[u.GetUserId()]; ok {
		err := activeUser.Conn.Close()
		if err != nil {
			return err
		}
	}
	c.ChatUsers[user.Id] = u

	err := SyncMissedMessages(u)
	if err != nil {
		return err
	}

	go u.ReadMessage()
	go u.WriteMessage()

	return nil
}

func SyncMissedMessages(user *ChatUser) error {
	messages, err := models.GetMissedMessages(&user.User)
	if err != nil {
		return err
	}

	for _, message := range messages {
		err := user.Conn.WriteMessage(websocket.TextMessage, []byte(message.Message))

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Chat) RemoveUser(user *ChatUser) {
	close(user.MessageChanel)
	delete(c.ChatUsers, user.GetUserId())

	repositories.UpdateUserLastActivity(&user.User)
}

func (c *Chat) GetActiveUsers() map[models.UserId]*ChatUser {
	return c.ChatUsers
}

func (cu ChatUser) ReadMessage() {
	defer func() {
		cu.Chat.Leavers <- cu
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

		message := &models.Message{Message: string(p)}
		_, err = message.Save()
		if err != nil {
			break
		}

		cu.Chat.MessageBroadcast <- p
	}
}

func (cu ChatUser) WriteMessage() {
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
