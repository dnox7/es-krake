package middleware

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"

	authEntities "github.com/dpe27/es-krake/internal/domain/auth/entity"
	authRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	authService "github.com/dpe27/es-krake/internal/domain/auth/service"
	enterpriseRepo "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	platformRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	domainerr "github.com/dpe27/es-krake/internal/domain/shared/errors"
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
	enterpriseAccountRepo  enterpriseRepo.EnterpriseAccountRepository
	permissionService      authService.PermissionService
	platformAccountRepo    platformRepo.PlatformAccountRepository
	roleRepo               authRepo.RoleRepository
}

func NewPermMiddleware(
	accessOpService authService.AccessOperationService,
	accessReqRepo authRepo.AccessRequirementRepository,
	entAccRepo enterpriseRepo.EnterpriseAccountRepository,
	pfAccRepo platformRepo.PlatformAccountRepository,
	permService authService.PermissionService,
	roleRepo authRepo.RoleRepository,
) *permMiddleware {
	return &permMiddleware{
		logger:                 log.With("middleware", "check_permission"),
		accessOperationService: accessOpService,
		accessRequirementRepo:  accessReqRepo,
		enterpriseAccountRepo:  entAccRepo,
		permissionService:      permService,
		platformAccountRepo:    pfAccRepo,
		roleRepo:               roleRepo,
	}
}

func (pm *permMiddleware) CheckPfPerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		pm.checkPerm(c, true)
	}
}

func (pm *permMiddleware) CheckEntPerm() gin.HandlerFunc {
	return func(c *gin.Context) {
		pm.checkPerm(c, false)
	}
}

type taskErr struct {
	id  int
	err error
}

func (pm *permMiddleware) checkPerm(c *gin.Context, isPfAcc bool) {
	var (
		wg          sync.WaitGroup
		permissions []authEntities.Permission
		operations  []authEntities.AccessOperation
	)
	errChan := make(chan taskErr, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		operations, err = pm.getReqOperationsFromReq(c)
		if err != nil {
			errChan <- taskErr{id: 1, err: err}
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		permissions, err = pm.getPermssionsFromKcUserID(c, isPfAcc)
		if err != nil {
			errChan <- taskErr{id: 2, err: err}
		}
	}()

	wg.Wait()
	close(errChan)
	if te, ok := <-errChan; ok {
		if te.id == 1 {
			nethttp.AbortWithInternalServerErrorResponse(c, te.err.Error(), nil, nil)
		} else {
			nethttp.AbortWithForbiddenResponse(c, te.err.Error(), nil, nil)
		}
		return
	}

	ok, err := authEntities.PermissionSlice(permissions).HasRequiredOperations(operations)
	if err != nil {
		nethttp.AbortWithInternalServerErrorResponse(c, err.Error(), nil, nil)
		return
	}

	if !ok {
		nethttp.AbortWithForbiddenResponse(c, http.StatusText(http.StatusForbidden), nil, nil)
	} else {
		c.Next()
	}
}

func (pm *permMiddleware) getPermssionsFromKcUserID(c *gin.Context, isPfAcc bool) ([]authEntities.Permission, error) {
	kcUserID, err := getKcUserIDFromReq(c)
	if err != nil {
		return nil, err
	}

	var (
		roleID     int
		roleTypeID int
	)
	if isPfAcc {
		pfAcc, err := pm.platformAccountRepo.TakeByConditions(c, map[string]interface{}{
			"kc_user_id": kcUserID,
		}, nil)
		if err != nil {
			return nil, err
		}

		roleID = pfAcc.RoleID
		roleTypeID = int(authRepo.PlatformRoleType)
	} else {
		entAcc, err := pm.enterpriseAccountRepo.TakeByConditions(c, map[string]interface{}{
			"kc_user_id": kcUserID,
		}, nil)
		if err != nil {
			return nil, err
		}

		roleID = entAcc.RoleID
		roleTypeID = int(authRepo.EnterpriseRoleType)
	}

	isPresent, err := pm.roleRepo.CheckExists(c, map[string]interface{}{
		"id":           roleID,
		"role_type_id": roleTypeID,
	}, nil)
	if err != nil {
		return nil, err
	}
	if !isPresent {
		return nil, domainerr.ErrorNotFound
	}

	return pm.permissionService.GetPermissionsWithRoleID(c, roleID)
}

func (pm *permMiddleware) getReqOperationsFromReq(c *gin.Context) ([]authEntities.AccessOperation, error) {
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
