package middleware

import (
	"errors"
	"strings"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak"
	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type authenMiddleware struct {
	logger       *log.Logger
	cfg          *config.Config
	kcKeyService keycloak.KcKeyService
}

func NewAuthenMiddleware(
	cfg *config.Config,
	kcKeyService keycloak.KcKeyService,
) *authenMiddleware {
	return &authenMiddleware{
		logger:       log.With("middleware", "check_authentication"),
		cfg:          cfg,
		kcKeyService: kcKeyService,
	}
}

func (am *authenMiddleware) CheckAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := getTokenFromReq(c)
		if tokenStr == "" {
			nethttp.AbortWithBadRequestResponse(c, map[string]interface{}{
				"access_token": utils.ErrorInputRequired,
			}, nil, nil)
			return
		}

		unauthorizedMsg := map[string]interface{}{
			"access_token": utils.ErrorInputFail,
		}

		token, err := jwt.DecodeJWTUnverified(tokenStr)
		if err != nil {
			nethttp.AbortWithUnauthorizedResponse(c, unauthorizedMsg, nil, nil)
			return
		}

		keyID := token.Header["kid"]
		if keyID == nil {
			nethttp.AbortWithUnauthorizedResponse(c, unauthorizedMsg, nil, nil)
			return
		}

		subject, err := token.Claims.GetSubject()
		if err != nil {
			nethttp.AbortWithUnauthorizedResponse(c, unauthorizedMsg, nil, nil)
			return
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			nethttp.AbortWithUnauthorizedResponse(c, unauthorizedMsg, nil, nil)
			return
		}

		pubKeys, ok := am.kcKeyService.GetJWKSKeysFromCache(issuer)
		_, err = am.kcKeyService.ExtractKey(keyID.(string), pubKeys)

		if !ok || err != nil {
			pubKeys, err := am.kcKeyService.GetJWKSKeysFromEndpoint(c, issuer)
			if err != nil {
				nethttp.AbortWithUnauthorizedResponse(c, unauthorizedMsg, nil, nil)
				return
			}
			am.kcKeyService.CacheJWKSKeyForRealm(issuer, pubKeys)
		}

		pubKey, err := am.kcKeyService.ExtractKey(keyID.(string), pubKeys)
		if err != nil {
			nethttp.AbortWithUnauthorizedResponse(c, err.Error(), nil, nil)
			return
		}

		claims, err := jwt.Verify(tokenStr, pubKey, am.cfg)
		if err != nil {
			if strings.Contains(err.Error(), jwt.ErrorTokenExpiredMsg) || errors.Is(err, jwtV5.ErrTokenExpired) {
				nethttp.AbortWithRequestTimeoutResponse(c, jwt.ErrorTokenExpiredMsg, nil, nil)
				return
			}

			nethttp.AbortWithUnauthorizedResponse(c, err.Error(), nil, nil)
			return
		}

		c.Set("kc_user_id", subject)
		c.Set("user", claims)
		c.Next()
	}
}

func getTokenFromReq(c *gin.Context) string {
	authHeader := c.Request.Header.Get(nethttp.HeaderAuthorization)
	if authHeader != "" && strings.Index(authHeader, nethttp.AuthSchemeBearer) == 0 {
		return authHeader[7:]
	}
	return ""
}
