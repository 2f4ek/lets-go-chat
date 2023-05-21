package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

type CustomLogger struct{}

func (l *CustomLogger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 599:
			{
				l.errorLogger(clientIP, method, path, statusCode)
			}
		default:
			l.defaultLogger(clientIP, method, path, statusCode)
		}

	}
}

func (l *CustomLogger) defaultLogger(clientIP string, method string, path string, statusCode int) {
	log.Printf("[Logger] %v | %v | %v | %v", clientIP, method, path, statusCode)
}

func (l *CustomLogger) errorLogger(clientIP string, method string, path string, statusCode int) {
	log.Printf("[ERROR] %v | %v | %v | %v", clientIP, method, path, statusCode)
}

func (l *CustomLogger) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[SERVER ERROR]")
				log.Printf("[PANIC] Error: %v", err)

				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
