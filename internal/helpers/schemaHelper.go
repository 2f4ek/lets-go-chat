package helpers

import (
	"github.com/gin-gonic/gin"
)

func GetSchema(ctx *gin.Context) string {
	scheme := "http://"
	if ctx.Request.TLS != nil {
		scheme = "https://"
	}

	return scheme
}
