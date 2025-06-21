package router

import (
	buyer "github.com/dpe27/es-krake/internal/interfaces/api/handler/buyer"
	"github.com/dpe27/es-krake/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

func BindBuyerRoute(
	router *gin.RouterGroup,
	handler *buyer.BuyerHandler,
	authenMiddleware *middleware.AuthenMiddleware,
	permMiddleware *middleware.PermMiddleware,
) {
	router.POST("/login", handler.PostBuyerLogin)
	router.POST("/signup", handler.PostBuyerSignup)
	router.POST("/auth/refresh", handler.PostRefreshTokenBuyer)

	router.Use(authenMiddleware.CheckAuthentication())
	{
		router.POST("/auth/logout", handler.PostLogout)
	}
}
