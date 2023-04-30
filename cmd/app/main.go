package main

import (
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user", handlers.RegisterUser)
	router.POST("/user/login", handlers.LoginUser)
	
	router.Run()
}
