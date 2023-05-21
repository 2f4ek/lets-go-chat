package middlewares

import "github.com/gin-gonic/gin"

type ILogger interface {
	Logger() gin.HandlerFunc
	Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc
	defaultLogger(clientIP string, method string, path string, statusCode int)
	errorLogger(clientIP string, method string, path string, statusCode int)
}
