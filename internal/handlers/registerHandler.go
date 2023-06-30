package handlers

import (
	"fmt"
	"net/http"

	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gin-gonic/gin"
)

var minPasswordLength = 8
var minLoginLength = 4
var RHInstance *RegisterHandler

type ICreateUserRequest interface {
	Validate() bool
}

type ICreateUserResponse interface {
	RegisterUser(ctx *gin.Context)
}

// CreateUserRequest
// @Description Create user request
type CreateUserRequest struct {
	// UserName of existing user
	UserName string `json:"userName"`
	// Password of existing user
	Password string `json:"password"`
}

// CreateUserResponse
// @Description Response with registered user data
type CreateUserResponse struct {
	// UserName name of new user
	UserName string `json:"userName"`
	// ID of new user
	ID models.UserId `json:"id"`
}

type RegisterHandler struct {
	ur *repositories.UserRepository
}

func ProvideRegisterHandler(ur *repositories.UserRepository) *RegisterHandler {
	once.Do(func() {
		RHInstance = &RegisterHandler{ur: ur}
	})
	return RHInstance
}

func (r *CreateUserRequest) Validate() bool {
	return len(r.UserName) > minLoginLength && len(r.Password) > minPasswordLength
}

// RegisterUser godoc
// @Summary Registration
// @Description Register user by userName and password
// @Schemes http https
// @Param request body handlers.CreateUserRequest true "query params"
// @failure 400 {string} string "Error message"
// @Success 200 {object} handlers.CreateUserResponse
// @Router /user [POST]
func (r *RegisterHandler) RegisterUser(ctx *gin.Context) {
	userRequest := &CreateUserRequest{}
	if err := ctx.Bind(userRequest); err != nil {
		ctx.String(http.StatusBadRequest, "Bad request, empty username or password")
		return
	}

	ok := userRequest.Validate()
	if !ok {
		ctx.String(http.StatusBadRequest,
			fmt.Sprintf(
				"username should contain more than %s chars and password should contain more than %s chars",
				string(rune(minLoginLength)), string(rune(minPasswordLength))))
		return
	}

	user, userExists := r.ur.CreateUser(userRequest.UserName, userRequest.Password)
	if userExists == true {
		ctx.String(http.StatusBadRequest, "Username already taken")
		return
	}

	r.ur.AppendUser(*user)

	ctx.JSON(201, CreateUserResponse{UserName: user.Name, ID: user.Id})
}
