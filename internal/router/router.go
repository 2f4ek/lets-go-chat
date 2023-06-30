package router

import (
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ch *handlers.ChatHandler
	rh *handlers.RegisterHandler
	ah *handlers.AuthHandler
}

func ProvideRouter(
	ch *handlers.ChatHandler,
	rh *handlers.RegisterHandler,
	ah *handlers.AuthHandler,
) Router {
	return Router{ch: ch, rh: rh, ah: ah}
}

func (r *Router) InitRoutes(app *gin.Engine) {
	app.GET("/ws", r.ch.WsInit)
	app.POST("/user", r.rh.RegisterUser)
	app.POST("/user/login", r.ah.LoginUser)
	app.GET("/active-users", r.ch.ActiveUsers)
}
