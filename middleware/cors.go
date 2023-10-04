package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AcceptCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Authorization, *")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusNoContent)
		return
	}

	c.Next()
}
