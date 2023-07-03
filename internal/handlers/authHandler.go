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

type IAuthHandler interface {
	LoginUser(ctx *gin.Context)
}

type AuthHandler struct {
	ur *repositories.UserRepository
}

type ILoginRequest interface {
	validate() bool
}

// LoginRequest
// @Description Login request
type LoginRequest struct {
	// UserName should be available in the system
	UserName string `json:"userName"`
	// Password should be more than 4 chars
	Password string `json:"password"`
}

// LoginResponse
// @Description Login response with websocket url
type LoginResponse struct {
	// Url contains path to connect to websocket
	Url string `json:"url"`
}

func ProvideAuthHandler(ur *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{ur: ur}
}

func (r *LoginRequest) validate() bool {
	return len(r.UserName) > 0 && len(r.Password) > 0
}

// LoginUser godoc
// @Summary Login
// @Description Login user by userName and password
// @Schemes http https
// @Param request body handlers.LoginRequest true "query params"
// @failure 400 {string} string "Error message"
// @Success 200 {object} handlers.LoginResponse
// @Router /user/login [POST]
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
