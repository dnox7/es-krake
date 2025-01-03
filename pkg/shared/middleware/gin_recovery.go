package middleware

import (
	"fmt"
	"net/http"
	"pech/es-krake/pkg/dto"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinCustomRecovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(
		func(c *gin.Context, recovered any) {
			var errMsg string

			if recovered != nil {
				errMsg = fmt.Sprintf("%v", recovered)
			}
			stackTrace := stack()

			ctxLogger := logger.WithContext(c)

			ctxLogger.WithFields(logrus.Fields{
				"error":       errMsg,
				"path":        c.Request.URL.Path,
				"method":      c.Request.Method,
				"client_ip":   c.ClientIP(),
				"stack_trace": stackTrace,
			}).Errorln("Panic occured")

			errDto := &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: fmt.Sprintf("Panic occured: %v", errMsg),
					Details: stackTrace,
				},
			}

			c.JSON(http.StatusInternalServerError, errDto)
		},
	)
}

// stack: is dynamically adjusts buffer size to capture the complete stack trace
// tabs are removed, and the stack trace is split into individual lines
func stack() (stackTrace []string) {
	for size := 1024; ; size *= 2 {
		stackBuff := make([]byte, size)
		length := runtime.Stack(stackBuff, false)
		if length < size {
			trimST := strings.ReplaceAll(string(stackBuff[:length]), "\t", "")
			return strings.Split(trimST, "\n")
		}
	}
}
