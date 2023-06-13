package main

import (
	"github.com/2f4ek/lets-go-chat/database"
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	loadEnv()
	loadDatabase()

	app := gin.New()

	initMiddlewares(app)
	chat := handlers.InitChat()
	go chat.RunChat()
	router.InitRoutes(app)

	app.Run()
}

func initMiddlewares(app *gin.Engine) {
	app.Use(middlewares.Logger())
	app.Use(middlewares.Recovery())
}

func loadDatabase() {
	database.Connect()
	err := database.Database.AutoMigrate(&models.Message{})
	if err != nil {
		return
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
