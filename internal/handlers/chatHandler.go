package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/2f4ek/lets-go-chat/internal/chatModels"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type IChatHandler interface {
	InitChat() *chatModels.Chat
	WsInit(c *gin.Context)
	ActiveUsers(c *gin.Context)
}

type ChatHandler struct {
	ur   *repositories.UserRepository
	chat *chatModels.Chat
}

var (
	CInstance *ChatHandler
	once      sync.Once
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func ProvideChatHandler(ur *repositories.UserRepository, chat *chatModels.Chat) *ChatHandler {
	once.Do(func() {
		CInstance = &ChatHandler{ur: ur, chat: chat}
	})
	return CInstance
}

func (ch *ChatHandler) InitChat() *chatModels.Chat {
	return ch.chat
}

func (ch *ChatHandler) WsInit(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, fmt.Sprint("Token is required"))
		return
	}

	user := ch.ur.GetUserByToken(token)
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

	err = ch.chat.LoginUserToChat(user, ws)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	ch.ur.RevokeToken(user)
}

func (ch *ChatHandler) ActiveUsers(c *gin.Context) {
	result := make([]models.UserId, 0)
	for _, user := range ch.chat.GetActiveUsers() {
		result = append(result, user.GetUserId())
	}

	c.JSON(http.StatusOK, result)
}
