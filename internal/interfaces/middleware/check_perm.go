package middleware

import (
	platformRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/gin-gonic/gin"
)

func CheckPfPerm(
	pfAccRepo platformRepo.PlatformAccountRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CheckEntPerm() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
