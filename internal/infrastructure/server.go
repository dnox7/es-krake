package infrastructure

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/interfaces/middleware"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
)

func NewGinRouter(cfg *config.Config) *gin.Engine {
	if cfg.App.Env != utils.ProdEnv {
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
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		s := <-signals
		log.Info(context.Background(), "Received signal '%v'. Shutting down...", "s", s.String())

		if err := callback(); err != nil {
			log.Error(context.Background(), "Error during graceful shutdown", "detail", err)
		}
	}()
}
