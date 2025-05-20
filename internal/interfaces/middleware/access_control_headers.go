package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultMaxAge = 86400
)

func AccessControlHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-store")
		c.Header("X-Content-Type-Options", "nosiff")

		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		if requestedHeaders := c.Request.Header.Get("Access-Control-Request-Headers"); requestedHeaders != "" {
			c.Header("Access-Control-Allow-Headers", requestedHeaders)
		} else {
			c.Header("Access-Control-Allow-Headers", "*")
		}

		if requestedMethod := c.Request.Header.Get("Access-Control-Request-Method"); requestedMethod != "" {
			c.Header("Access-Control-Allow-Method", requestedMethod)
		} else {
			c.Header("Access-Control-Allow-Method", "*")
		}

		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Max-Age", strconv.Itoa(defaultMaxAge))
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}
