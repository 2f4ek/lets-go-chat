package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/pkg/hasher"
	"github.com/gin-gonic/gin"
)

var AHInstance *AuthHandler

type IAuthHandler interface {
	LoginUser(ctx *gin.Context)
}

type AuthHandler struct {
	ur *repositories.UserRepository
}

type ILoginRequest interface {
	validate() bool
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Url string `json:"url"`
}

func ProvideAuthHandler(ur *repositories.UserRepository) *AuthHandler {
	once.Do(func() {
		AHInstance = &AuthHandler{ur: ur}
	})
	return AHInstance
}

func (r *LoginRequest) validate() bool {
	return len(r.UserName) > 0 && len(r.Password) > 0
}

func (ah *AuthHandler) LoginUser(ctx *gin.Context) {
	loginRequest := &LoginRequest{}
	if err := ctx.Bind(loginRequest); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid username/password")
		return
	}

	if ok := loginRequest.validate(); !ok {
		ctx.String(http.StatusBadRequest, fmt.Sprint("Bad request, empty user name or password"))
		return
	}

	user, userExists := ah.ur.GetUser(loginRequest.UserName)
	if !userExists {
		ctx.String(http.StatusBadRequest, fmt.Sprint("User not founded"))
		return
	}

	if !hasher.CheckPasswordHash(loginRequest.Password, user.Password) {
		ctx.String(http.StatusBadRequest, fmt.Sprint("Wrong password"))
		return
	}

	ctx.Header("X-Rate-Limit", "999999")
	ctx.Header("X-Expires-After", time.Now().Add(time.Hour*1).UTC().String())

	token := helpers.GenerateSecureToken()
	ah.ur.UpdateToken(user, token)

	ctx.JSON(http.StatusOK,
		LoginResponse{Url: "wss://" + ctx.Request.Host + "/ws?token=" + token})
}
