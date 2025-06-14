package middleware

import (
	"fmt"
	"time"

	"github.com/dpe27/es-krake/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDKey = "request-id"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		upstreamRequestID := c.GetHeader(requestIDKey)
		if upstreamRequestID != "" {
			requestID = upstreamRequestID
		}

		// Add request-id to all servers
		ctx := log.AddLogValToCtx(c.Request.Context(), requestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set(requestIDKey, requestID)

		c.Next()
		duration := time.Since(start)

		if c.Request.RequestURI == "/health-check" && c.Writer.Status() < 400 {
			return
		}

		fields := map[string]interface{}{
			"host":      c.Request.Host,
			"duration":  duration.String(),
			"clientIp":  c.ClientIP(),
			"method":    c.Request.Method,
			"url":       c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"referer":   c.Request.Referer(),
			"userAgent": c.Request.UserAgent(),
		}

		// en: Add additional log context if available.
		mergeLogContext(fields, c)

		// en: Convert fields to key-value arguments for logging.
		args := make([]any, 0, len(fields)*2)
		for k, v := range fields {
			args = append(args, k, v)
		}

		log.With(args...).Info(c, "")
	}
}

func mergeLogContext(fields map[string]interface{}, c *gin.Context) {
	if logContext, exists := c.Keys["logContext"]; exists {
		if logContextMap, ok := logContext.(map[string]interface{}); ok {
			for key, value := range logContextMap {
				fields[key] = value
			}
		} else {
			log.Error(
				c, fmt.Sprintf(
					"Cannot include the logContext '%v' because it is not a map[string]interface{}.",
					logContext,
				),
			)
		}
	}
}
