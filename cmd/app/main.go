package main

import (
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()

	initMiddlewares(app)
	initRoutes(app)

	app.Run()
}

func initMiddlewares(app *gin.Engine) {
	app.Use(middlewares.Logger())
	app.Use(middlewares.Recovery())
}

func initRoutes(app *gin.Engine) {
	app.GET("/ws", handlers.WsInit)
	app.POST("/user", handlers.RegisterUser)
	app.POST("/user/login", handlers.LoginUser)
}
