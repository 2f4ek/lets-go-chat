package main

import (
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	logger := &middlewares.CustomLogger{}
	app.Use(logger.Logger())
	app.Use(logger.Recovery())

	app.POST("/user", handlers.RegisterUser)
	app.POST("/user/login", handlers.LoginUser)

	app.Run()
}
