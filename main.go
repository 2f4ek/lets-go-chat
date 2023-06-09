package main

import (
	"log"

	"github.com/2f4ek/lets-go-chat/database"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           2f4ek Lets Go Chat openAPI documentation
// @version         1.0
// @description     Chat.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	loadEnv()
	runMigrations(database.Database{})

	routes, err := InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	app := gin.New()
	initMiddlewares(app)
	routes.InitRoutes(app)
	app.Run()
}

func initMiddlewares(app *gin.Engine) {
	app.Use(middlewares.Logger())
	app.Use(middlewares.Recovery())
}

func runMigrations(db database.Database) {
	db.Connect()
	err := db.Database.AutoMigrate(&models.Message{})
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
