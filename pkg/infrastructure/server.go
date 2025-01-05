package infrastructure

import (
	"os"
	"pech/es-krake/pkg/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB, logger *logrus.Logger) *gin.Engine {
	if os.Getenv("PE_DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.Recovery(logger),
		middleware.AccessControlHeaders(),
		middleware.GinCustomRecovery(logger),
	)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	return router
}
