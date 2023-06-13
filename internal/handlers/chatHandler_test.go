package handlers_test

import (
	"github.com/2f4ek/lets-go-chat/internal/chatModels"
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	userRepo "github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitChat(t *testing.T) {
	chat := handlers.InitChat()

	assert.Same(t, chat, handlers.InitChat())
	assert.IsType(t, &chatModels.Chat{}, chat)
}

func TestActiveUsers(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/active-users", nil)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWsInitUserTokenError(t *testing.T) {
	user, _ := userRepo.CreateUser("TestWsInitUserTokenError", "password1")
	token := helpers.GenerateSecureToken()
	user.Token = token
	userRepo.AppendUser(*user)

	request := httptest.NewRequest(http.MethodGet, "/ws", nil)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWsInitServerError(t *testing.T) {
	user, _ := userRepo.CreateUser("TestWsInitUserExistsError", "password1")
	token := helpers.GenerateSecureToken()
	user.Token = token
	userRepo.AppendUser(*user)

	request := httptest.NewRequest(http.MethodGet, "/ws", nil)
	request.Header.Add("Content-Type", "application/json")
	query := request.URL.Query()
	query.Add("token", token)
	app := gin.New()
	router.InitRoutes(app)

	request.URL.RawQuery = query.Encode()
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWsInitUserNotExistsError(t *testing.T) {
	token := helpers.GenerateSecureToken()
	request := httptest.NewRequest(http.MethodGet, "/ws", nil)
	query := request.URL.Query()
	query.Add("token", token)
	app := gin.New()
	router.InitRoutes(app)

	request.URL.RawQuery = query.Encode()
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)
	t.Logf("%v", w.Result().Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
