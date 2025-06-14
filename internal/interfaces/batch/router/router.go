package router

import (
	"github.com/dpe27/es-krake/internal/interfaces/batch/handler"
	"github.com/gin-gonic/gin"
)

func BindBatchRoutes(
	router *gin.RouterGroup,
	h *handler.BatchHander,
) {
	router.POST("/batch", h.Batch)
	router.POST("/dead-letter", h.DeadLetter)
}
