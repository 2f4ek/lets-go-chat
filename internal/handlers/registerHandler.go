package handlers

import (
	"fmt"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

var minPasswordLength = 6
var minLoginLength = 4

type CreateUserRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

func (r *CreateUserRequest) Validate() bool {
	return len(r.UserName) > minLoginLength && len(r.Password) > minPasswordLength
}

func RegisterUser(ctx *gin.Context) {
	userRequest := &CreateUserRequest{}
	if err := ctx.Bind(userRequest); err != nil {
		ctx.String(http.StatusBadRequest, "Bad request")
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

	user, userExists := repositories.CreateUser(userRequest.UserName, userRequest.Password)
	if userExists == true {
		ctx.String(http.StatusBadRequest, "User name already taken")
		return
	}

	repositories.AppendUser(*user)

	ctx.JSON(201, CreateUserResponse{UserName: user.Name, Token: user.Token})
}
