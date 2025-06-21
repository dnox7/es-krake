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

type PlatformAuthToken struct {
	Token       kcdto.TokenEndpointResp
	RealmName   string
	Permissions []entity.Permission
}

func (u *AuthUsecase) LoginPlatform(ctx context.Context, username, password string) (*PlatformAuthToken, error) {
	platformClient := u.deps.KcClientService.GetPlatformClient()
	realmName := platformClient["realm_name"]
	clientID := platformClient["client_id"]

	tokenResp, err := u.deps.KcTokenService.GetTokenWithPassword(
		ctx,
		realmName,
		clientID,
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

	return &PlatformAuthToken{
		Token:       tokenResp,
		RealmName:   realmName,
		Permissions: permissions,
	}, nil
}
