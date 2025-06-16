package usecase

import (
	"context"
	"strings"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak/utils"
	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

const ErrorPlatformLoginFailed = "failed to login platform"

type LoginPlatformResponse struct {
	Token       kcdto.TokenEndpointResp
	Permissions []entity.Permission
}

func (u *AuthUsecase) LoginPlatform(ctx context.Context, username, password string) (*LoginPlatformResponse, error) {
	platformClient := u.deps.KcClientService.GetPlatformClient()

	tokenResp, err := u.deps.KcTokenService.GetTokenWithPassword(
		ctx,
		platformClient["realm_name"].(string),
		platformClient["client_id"].(string),
		username,
		password,
	)
	if err != nil {
		if strings.Contains(err.Error(), utils.KeycloakLoginInvalidCredential) {
			return nil, wraperror.NewValidationError(
				map[string]interface{}{
					"credentials": ErrorPlatformLoginFailed,
				},
				err,
			)
		}
		return nil, err
	}

	token, err := jwt.DecodeJWTUnverified(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	kcUserID := jwt.GetKeycloakUserID(token.Claims)
	platformAccount, err := u.deps.PlatformAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"kc_user_id": kcUserID,
	}, nil)
	if err != nil {
		return nil, err
	}

	permissions, err := u.deps.PermissionService.GetPermissionsWithRoleID(ctx, platformAccount.RoleID)
	if err != nil {
		return nil, err
	}

	return &LoginPlatformResponse{
		Token:       tokenResp,
		Permissions: permissions,
	}, nil
}
