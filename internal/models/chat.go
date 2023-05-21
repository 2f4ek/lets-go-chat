package models

import "github.com/gorilla/websocket"

type Chat struct {
	ChatUsers map[string]ChatUser
}

func (chat *Chat) AddUserToChat(user User, conn *websocket.Conn) {
	chat.ChatUsers[user.Token] = ChatUser{
		Token:    user.Token,
		Conn:     conn,
		UserName: user.Name,
		UserId:   user.Id,
	}
}

func (chat *Chat) RemoveUser(token string) {
	delete(chat.ChatUsers, token)
}

func (chat *Chat) GetActiveUsers() map[string]ChatUser {
	return chat.ChatUsers
}
