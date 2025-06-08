package middleware

import (
	"net/http"
	"strconv"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/gin-gonic/gin"
)

const (
	defaultMaxAge = 86400
)

func AccessControlHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(nethttp.HeaderAccessControlAllowCredentials, "true")
		c.Header(nethttp.HeaderCacheControl, "no-store")
		c.Header(nethttp.HeaderPragma, "no-store")
		c.Header(nethttp.HeaderXContentTypeOptions, "nosiff")

		if origin := c.Request.Header.Get(nethttp.HeaderOrigin); origin != "" {
			c.Header(nethttp.HeaderAccessControlAllowOrigin, origin)
		} else {
			c.Header(nethttp.HeaderAccessControlAllowOrigin, "*")
		}

		if requestedHeaders := c.Request.Header.Get(nethttp.HeaderAccessControlRequestHeaders); requestedHeaders != "" {
			c.Header(nethttp.HeaderAccessControlAllowHeaders, requestedHeaders)
		} else {
			c.Header(nethttp.HeaderAccessControlAllowHeaders, "*")
		}

		if requestedMethod := c.Request.Header.Get(nethttp.HeaderAccessControlRequestMethod); requestedMethod != "" {
			c.Header(nethttp.HeaderAccessControlAllowMethod, requestedMethod)
		} else {
			c.Header(nethttp.HeaderAccessControlAllowMethod, "*")
		}

		if c.Request.Method == http.MethodOptions {
			c.Header(nethttp.HeaderAccessControlMaxAge, strconv.Itoa(defaultMaxAge))
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}
