package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"pech/es-krake/pkg/log"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.
					With("stack", string(debug.Stack())).
					Error(c.Request.Context(), "panic was trigger", "error", err)

				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") || strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					c.Abort()
				}

				c.JSON(http.StatusInternalServerError, "Internal Error")
				return
			}
		}()
		c.Next()
	}
}
