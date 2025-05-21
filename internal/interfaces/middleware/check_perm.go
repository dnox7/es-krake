package middleware

import (
	"errors"
	"regexp"
	"strings"

	authEntities "github.com/dpe27/es-krake/internal/domain/auth/entity"
	authRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	authService "github.com/dpe27/es-krake/internal/domain/auth/service"
	platformRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	ErrGetAccessReqCode      = "failed to take access_requirement_code from request"
	ErrAccessReqCodeNotFound = "access_requirement_code doesn not exists"
)

type permMiddleware struct {
	logger                 *log.Logger
	accessOperationService authService.AccessOperationService
	accessRequirementRepo  authRepo.AccessRequirementRepository
	platformAccountRepo    platformRepo.PlatformAccountRepository
	roleRepo               authRepo.RoleRepository
}

func NewPermMiddleware(
	accessOpService authService.AccessOperationService,
	accessReqRepo authRepo.AccessRequirementRepository,
	pfAccRepo platformRepo.PlatformAccountRepository,
	roleRepo authRepo.RoleRepository,
) *permMiddleware {
	return &permMiddleware{
		logger:                 log.With("middleware", "check_permission"),
		accessOperationService: accessOpService,
		platformAccountRepo:    pfAccRepo,
		roleRepo:               roleRepo,
	}
}

func (pm *permMiddleware) CheckPfPerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		kcUserID, err := getKcUserIDFromReq(c)
		if err != nil {
			nethttp.AbortWithForbiddenResponse(c, err.Error(), nil, nil)
		}

		acc, err := pm.platformAccountRepo.TakeByConditions(c, map[string]interface{}{
			"kc_user_id": kcUserID,
		}, nil)
		if err != nil {
			nethttp.AbortWithForbiddenResponse(c, err.Error(), nil, nil)
		}

		exists, err := pm.roleRepo.CheckExists(c, map[string]interface{}{
			"id": acc.RoleID,
		}, nil)
		if !exists {
			nethttp.AbortWithForbiddenResponse(c, err.Error(), nil, nil)
		}

	}
}

func (pm *permMiddleware) CheckEntPerm() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (pm *permMiddleware) getRequirementOperationsFromReq(c *gin.Context) ([]authEntities.AccessOperation, error) {
	var accessReqCode string
	for key, code := range routeMatcherIndex {
		pattern, err := parsePatternKey(key)
		if err != nil {
			return nil, err
		}

		checkMethodMatch, err := regexp.MatchString(pattern.Method, c.Request.Method)
		if err != nil {
			return nil, err
		}

		checkPathMatch, err := utils.CheckKeyMatch(pattern.Path, c.Request.URL.Path)
		if err != nil {
			return nil, err
		}

		if checkMethodMatch && checkPathMatch {
			accessReqCode = code
			break
		}
	}
	if accessReqCode == "" {
		pm.logger.Error(c, ErrGetAccessReqCode)
		return nil, errors.New(ErrGetAccessReqCode)
	}

	exists, err := pm.accessRequirementRepo.CheckExists(c, map[string]interface{}{
		"code": accessReqCode,
	}, nil)
	if err != nil {
		return nil, err
	}
	if !exists {
		pm.logger.Error(c, ErrAccessReqCodeNotFound, "code", accessReqCode)
		return nil, errors.New(ErrAccessReqCodeNotFound)
	}

	return pm.accessOperationService.GetOperationsWithAccessReqCode(c, accessReqCode)
}

func getKcUserIDFromReq(c *gin.Context) (string, error) {
	tokenStr := getTokenFromReq(c)
	if tokenStr == "" {
		return "", jwt.ErrMissingToken
	}

	token, err := jwt.DecodeJWTUnverified(tokenStr)
	if err != nil {
		return "", jwt.ErrInvalidToken
	}

	kcUserID := jwt.GetKeycloakUserID(token.Claims)
	if kcUserID == "" {
		return "", jwt.ErrInvalidToken
	}

	return kcUserID, nil
}

func getTokenFromReq(c *gin.Context) string {
	authHeader := c.Request.Header.Get(nethttp.HeaderAuthorization)
	if authHeader != "" && strings.Index(authHeader, nethttp.AuthSchemeBearer) == 0 {
		return authHeader[7:]
	}
	return ""
}
