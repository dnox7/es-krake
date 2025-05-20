package infrastructure

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dpe27/es-krake/internal/interfaces/middleware"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	if os.Getenv("PE_DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.Recovery(),
		middleware.AccessControlHeaders(),
		middleware.GinCustomRecovery(),
	)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	return router
}

func OnShutdown(callback func() error) {
	go (func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		s := <-signals
		slog.InfoContext(context.Background(), "Received signal '%v'. Shutting down...", "s", s.String())

		if err := callback(); err != nil {
			slog.ErrorContext(context.Background(), "Error during graceful shutdown", "detail", err)
		}
	})()
}
