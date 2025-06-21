package router

import (
	ent "github.com/dpe27/es-krake/internal/interfaces/api/handler/enterprise"
	"github.com/dpe27/es-krake/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

func BindEnterpriseRoute(
	router *gin.RouterGroup,
	handler *ent.EnterpriseHandler,
	authenMiddleware *middleware.AuthenMiddleware,
	permMiddleware *middleware.PermMiddleware,
) {
	router.POST("/:enterprise_id/login", handler.PostLogin)
	router.POST("/:enterprise_id/auth/refresh", handler.PostRefreshTokenEnterprise)

	router.Use(authenMiddleware.CheckAuthentication())
	{
		router.POST("/:enterprise_id/auth/logout", handler.PostLogout)
	}
}
