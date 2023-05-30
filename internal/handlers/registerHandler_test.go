package handlers_test

import (
	"bytes"
	userRepo "github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user",
		bytes.NewReader([]byte(`{"userName": "userName1", "password": "password11"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegisterUserMissingFieldError(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user",
		bytes.NewReader([]byte(`[]`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUserValidationError(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/user",
		bytes.NewReader([]byte(`{"userName": "userName", "password": "pass"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUserExistingUserError(t *testing.T) {
	user, _ := userRepo.CreateUser("registeredUser", "password1")
	userRepo.AppendUser(*user)

	request := httptest.NewRequest(
		http.MethodPost,
		"/user",
		bytes.NewReader([]byte(`{"userName": "registeredUser", "password": "password1"}`)),
	)
	request.Header.Add("Content-Type", "application/json")

	app := gin.New()
	router.InitRoutes(app)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
