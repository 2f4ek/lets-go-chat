package logger

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Log(c *gin.Context, message string, statusCode int) {
	createResponseMessage(c, message, statusCode)
}

func createResponseMessage(c *gin.Context, message string, statusCode int) {
	clientIP := c.ClientIP()
	method := c.Request.Method
	path := c.Request.URL.Path

	log.Printf("%v %v | %v | %v | %v", message, clientIP, method, path, statusCode)
}
