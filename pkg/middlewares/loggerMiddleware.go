package middlewares

import (
	"github.com/2f4ek/lets-go-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"log"
)

const errorLog = "[ERROR]"
const defaultLog = "[LOGGER]"

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		statusCode := c.Writer.Status()

		switch {
		case statusCode >= 400 && statusCode <= 599:
			{
				logger.Log(c, errorLog, statusCode)
			}
		default:
			logger.Log(c, defaultLog, statusCode)
		}
	}
}

func Recovery() gin.HandlerFunc {
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
