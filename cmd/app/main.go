package main

import (
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()

	initMiddlewares(app)
	router.InitRoutes(app)

	app.Run()
}

func initMiddlewares(app *gin.Engine) {
	app.Use(middlewares.Logger())
	app.Use(middlewares.Recovery())
}
