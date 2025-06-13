package router

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/handler/pf"
	"github.com/gin-gonic/gin"
)

func BindPlatformRoute(router *gin.RouterGroup, handler *pf.PlatformHandler) {
	router.GET("/role_type/:role_type_id", handler.GetRoleType)
}
