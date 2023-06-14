package handlers

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/internal/chatModels"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	chat     *chatModels.Chat
	once     sync.Once
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func InitChat() *chatModels.Chat {
	once.Do(func() {
		chat = &chatModels.Chat{
			MessageBroadcast: make(chan []byte),
			Logout:           make(chan chatModels.ChatUser),
			Login:            make(chan chatModels.LoginUser),
			ChatUsers:        make(map[models.UserId]*chatModels.ChatUser),
		}
	})

	return chat
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

	err = chat.LoginUserToChat(user, ws)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	repositories.RevokeToken(user)
}

func ActiveUsers(c *gin.Context) {
	result := make([]models.UserId, 0)
	for _, user := range chat.GetActiveUsers() {
		result = append(result, user.GetUserId())
	}

	c.JSON(http.StatusOK, result)
}
