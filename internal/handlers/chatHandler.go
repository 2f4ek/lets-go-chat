package handlers

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var chat *models.Chat
var once sync.Once

func InitChat() *models.Chat {
	once.Do(func() {
		chat = &models.Chat{
			ChatUsers: make(map[string]models.ChatUser),
		}
	})

	return chat
}

func reader(conn *websocket.Conn, user *models.User) {
	messageType, p, err := conn.ReadMessage()
	chat := InitChat()
	if err != nil {
		log.Println(err)
		chat.RemoveUser(user.Token)
		return
	}

	for _, chat := range chat.GetActiveUsers() {
		if err := chat.Conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func WsInit(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, fmt.Sprint("Token is required"))
		return
	}

	user := repositories.GetUserByToken(token)
	if user == nil {
		c.String(http.StatusBadRequest, fmt.Sprint("Token is invalid"))
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	chat := InitChat()
	chat.AddUserToChat(*user, ws)
	repositories.RevokeToken(user)

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			logger.Log(c, "WebSocket error", http.StatusInternalServerError)
		}
	}(ws)

	for {
		reader(ws, user)
	}
}

func ActiveUsers(c *gin.Context) {
	chat := InitChat()
	users := chat.ChatUsers
	c.JSON(http.StatusOK, users)
}
