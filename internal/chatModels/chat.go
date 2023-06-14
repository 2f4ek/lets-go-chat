package chatModels

import (
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gorilla/websocket"
)

type Chat struct {
	ChatUsers        map[models.UserId]*ChatUser
	Logout           chan ChatUser
	Login            chan LoginUser
	MessageBroadcast chan []byte
}

type LoginUser struct {
	User *models.User
	Conn *websocket.Conn
}

func (c *Chat) RunChat() {
	for {
		select {
		case user := <-c.Logout:
			if _, ok := c.ChatUsers[user.GetUserId()]; ok {
				c.RemoveUser(&user)
			}
		case user := <-c.Login:
			c.AddUserToChat(&user)
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

func (c *Chat) AddUserToChat(loginUser *LoginUser) {
	u := &ChatUser{
		Conn:          loginUser.Conn,
		Chat:          c,
		MessageChanel: make(chan []byte),
		User:          *loginUser.User,
	}
	if activeUser, ok := c.ChatUsers[u.GetUserId()]; ok {
		err := activeUser.Conn.Close()
		if err != nil {
			return
		}
	}

	c.ChatUsers[u.GetUserId()] = u
	go u.SyncMissedMessages()
	go u.ReadMessage()
	go u.WriteMessage()
}

func (c *Chat) LoginUserToChat(user *models.User, conn *websocket.Conn) error {
	if activeUser, ok := c.ChatUsers[user.Id]; ok {
		err := activeUser.Conn.Close()
		if err != nil {
			return err
		}
	}

	loginUser := LoginUser{User: user, Conn: conn}
	c.Login <- loginUser

	return nil
}

func (cu ChatUser) SyncMissedMessages() {
	messages, err := models.GetMissedMessages(&cu.User)
	if err != nil {
		return
	}

	for _, message := range messages {
		err := cu.Conn.WriteMessage(websocket.TextMessage, []byte(message.Message))

		if err != nil {
			return
		}
	}

	return
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
		cu.Chat.Logout <- cu
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
