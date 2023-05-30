package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	userRepo "github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginUser(t *testing.T) {
	user, _ := userRepo.CreateUser("newUser", "password1")
	userRepo.AppendUser(*user)

	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`{"userName": "newUser", "password": "password1"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	body, _ := io.ReadAll(w.Body)

	url := &handlers.LoginResponse{}
	err := json.Unmarshal(body, url)

	assert.Nil(t, err)
	assert.NotEmpty(t, url.Url)
}

func TestLoginUserValidation(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`[]`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUserValidationError(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`{"userName": "", "password": "password1"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUserNotFoundError(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`{"userName": "notExistingUser", "password": "password1"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUserWrongPasswordError(t *testing.T) {
	user, _ := userRepo.CreateUser("existingUser", "password1")
	userRepo.AppendUser(*user)

	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`{"userName": "existingUser", "password": "wrongPassword"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
