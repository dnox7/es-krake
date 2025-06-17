package router

import (
	pf "github.com/dpe27/es-krake/internal/interfaces/api/handler/platform"
	"github.com/dpe27/es-krake/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

func BindPlatformRoute(
	router *gin.RouterGroup,
	handler *pf.PlatformHandler,
	authenMiddleware *middleware.AuthenMiddleware,
	permMiddleware *middleware.PermMiddleware,
) {
	router.GET("/role_type/:role_type_id", handler.GetRoleType)
	router.POST("/login", handler.PostPlatformLogin)

	router.Use(authenMiddleware.CheckAuthentication())
	{
		router.POST("/auth/logout", handler.PostLogout)
	}
}
