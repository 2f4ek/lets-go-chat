package router

import (
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	app.GET("/ws", handlers.WsInit)
	app.POST("/user", handlers.RegisterUser)
	app.POST("/user/login", handlers.LoginUser)
	app.GET("/active-users", handlers.ActiveUsers)
}
