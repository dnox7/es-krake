package infrastructure

import (
	"os"
	"pech/es-krake/pkg/middleware"

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
