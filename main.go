package main

import (
	"github.com/2f4ek/lets-go-chat/database"
	"github.com/2f4ek/lets-go-chat/internal/models"
	"github.com/2f4ek/lets-go-chat/pkg/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
)

// @title           2f4ek Lets Go Chat openAPI documentation
// @version         1.0
// @description     Chat.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	CreateTrace()
	RunChat()
}

func CreateTrace() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()
}

func RunChat() {
	loadEnv()
	runMigrations(database.Database{})

	routes, err := InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	app := gin.New()
	pprof.Register(app)

	initMiddlewares(app)
	routes.InitRoutes(app)

	err = app.Run()
	if err != nil {
		panic(err)
	}
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
