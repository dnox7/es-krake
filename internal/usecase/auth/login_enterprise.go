package usecase

import (
	"context"
	"strings"

	authEntities "github.com/dpe27/es-krake/internal/domain/auth/entity"
	enterpriseEntities "github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak/utils"
	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

const ErrorLoginEnterpriseFailed = "failed to login enterprise"

type LoginEnterpriseResponse struct {
	Token             kcdto.TokenEndpointResp
	EnterpriseAccount *enterpriseEntities.EnterpriseAccount
	RealmName         string
	Permissions       []authEntities.Permission
}

func (u *AuthUsecase) LoginEnterprise(ctx context.Context, enterpriseID int, username, password string) (*LoginEnterpriseResponse, error) {
	enterprise, err := u.deps.EnterpriseRepo.TakeByConditions(ctx, map[string]interface{}{
		"id":        enterpriseID,
		"is_active": true,
	}, nil)
	if err != nil {
		return nil, err
	}

	platformClient := u.deps.KcClientService.GetPlatformClient()

	platformTokenResp, err := u.deps.KcTokenService.GetTokenWithPassword(
		ctx,
		platformClient["realm_name"].(string),
		platformClient["client_id"].(string),
		username,
		password,
	)
	if err != nil && !strings.Contains(err.Error(), utils.KeycloakLoginInvalidCredential) {
		return nil, err
	}

	// is platform user
	if err == nil {
		platformToken, err := jwt.DecodeJWTUnverified(platformTokenResp.AccessToken)
		if err != nil {
			return nil, err
		}

		kcUserID := jwt.GetKeycloakUserID(platformToken.Claims)
		platformAccount, err := u.deps.PlatformAccountRepo.TakeByConditions(ctx, map[string]interface{}{
			"kc_user_id": kcUserID,
		}, nil)
		if err != nil {
			return nil, err
		}

		_, err = u.deps.PlatformAccountEnterpriseAccessRepo.TakeByConditions(ctx, map[string]interface{}{
			"platform_account_id": platformAccount.ID,
			"enabled":             true,
			"enterprise_id":       enterprise.ID,
		}, nil)
		if err != nil {
			return nil, err
		}

		permissions, err := u.deps.PermissionService.GetPermissionsWithRoleID(ctx, platformAccount.RoleID)
		if err != nil {
			return nil, err
		}

		return &LoginEnterpriseResponse{
			Token:             platformTokenResp,
			EnterpriseAccount: nil,
			RealmName:         platformClient["realm_name"].(string),
			Permissions:       permissions,
		}, nil
	}

	enterpriseTokenResp, err := u.deps.KcTokenService.GetTokenWithPassword(
		ctx,
		enterprise.KcRealmName,
		enterprise.KcClientID,
		username,
		password,
	)
	if err != nil {
		if strings.Contains(err.Error(), utils.KeycloakLoginInvalidCredential) ||
			strings.Contains(err.Error(), utils.KeycloakUnverifiedEmail) {
			return nil, wraperror.NewValidationError(
				map[string]interface{}{
					"credential": ErrorLoginEnterpriseFailed,
				},
				err,
			)
		}
		return nil, err
	}

	enterpriseToken, err := jwt.DecodeJWTUnverified(enterpriseTokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	kcUserID := jwt.GetKeycloakUserID(enterpriseToken.Claims)
	enterpriseAccount, err := u.deps.EnterpriseAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"enterprise_id": enterprise.ID,
		"kc_user_id":    kcUserID,
	}, nil)
	if err != nil {
		return nil, err
	}

	permissions, err := u.deps.PermissionService.GetPermissionsWithRoleID(ctx, enterpriseAccount.RoleID)
	if err != nil {
		return nil, err
	}

	return &LoginEnterpriseResponse{
		Token:             enterpriseTokenResp,
		EnterpriseAccount: &enterpriseAccount,
		RealmName:         enterprise.KcRealmName,
		Permissions:       permissions,
	}, nil
}
