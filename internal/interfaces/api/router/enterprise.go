package router

import (
	ent "github.com/dpe27/es-krake/internal/interfaces/api/handler/enterprise"
	"github.com/gin-gonic/gin"
)

func BindEnterpriseRoute(router *gin.RouterGroup, handler *ent.EnterpriseHandler) {
	router.POST("/:enterprise_id/login", handler.PostLogin)
}
