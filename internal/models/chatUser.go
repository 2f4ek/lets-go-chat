package models

import (
	"github.com/gorilla/websocket"
)

type ChatUser struct {
	Conn  *websocket.Conn
	Token string
}
