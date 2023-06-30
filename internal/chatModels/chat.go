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
	UserRepository   *repositories.UserRepository
	ChatMessageRep   *repositories.ChatMessageRepository
}

type LoginUser struct {
	User *models.User
	Conn *websocket.Conn
}

func ProvideChat(ur *repositories.UserRepository, cmr *repositories.ChatMessageRepository) *Chat {
	return &Chat{
		MessageBroadcast: make(chan []byte),
		Logout:           make(chan ChatUser),
		Login:            make(chan LoginUser),
		ChatUsers:        make(map[models.UserId]*ChatUser),
		UserRepository:   ur,
		ChatMessageRep:   cmr,
	}
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

func (c *Chat) RemoveUser(user *ChatUser) {
	close(user.MessageChanel)
	delete(c.ChatUsers, user.GetUserId())

	c.UserRepository.UpdateUserLastActivity(&user.User)
}

func (c *Chat) GetActiveUsers() map[models.UserId]*ChatUser {
	return c.ChatUsers
}
