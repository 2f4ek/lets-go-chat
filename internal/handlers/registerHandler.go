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
var RHInstanse *RegisterHandler

type ICreateUserRequest interface {
	Validate() bool
}

type ICreateUserResponse interface {
	RegisterUser(ctx *gin.Context)
}

type CreateUserRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserName string        `json:"userName"`
	Id       models.UserId `json:"id"`
}

type RegisterHandler struct {
	ur *repositories.UserRepository
}

func ProvideRegisterHandler(ur *repositories.UserRepository) *RegisterHandler {
	once.Do(func() {
		RHInstanse = &RegisterHandler{ur: ur}
	})
	return RHInstanse
}

func (r *CreateUserRequest) Validate() bool {
	return len(r.UserName) > minLoginLength && len(r.Password) > minPasswordLength
}

func (r *RegisterHandler) RegisterUser(ctx *gin.Context) {
	userRequest := &CreateUserRequest{}
	if err := ctx.Bind(userRequest); err != nil {
		ctx.String(http.StatusBadRequest, "Bad request, empty username or id")
		return
	}

	ok := userRequest.Validate()
	if !ok {
		ctx.String(http.StatusBadRequest,
			fmt.Sprintf(
				"user name should contain more than %s chars and password should contain more than %s chars",
				string(rune(minLoginLength)), string(rune(minPasswordLength))))
		return
	}

	user, userExists := r.ur.CreateUser(userRequest.UserName, userRequest.Password)
	if userExists == true {
		ctx.String(http.StatusBadRequest, "User name already taken")
		return
	}

	r.ur.AppendUser(*user)

	ctx.JSON(201, CreateUserResponse{UserName: user.Name, Id: user.Id})
}
