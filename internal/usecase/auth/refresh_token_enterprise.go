package usecase

import (
	"context"
	"errors"

	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) RefreshTokenEnterprise(ctx context.Context, enterpriseID int, cookies map[string]interface{}) (*EnterpriseAuthToken, error) {
	enterprise, err := u.deps.EnterpriseRepo.TakeByConditions(ctx, map[string]interface{}{
		"id":        enterpriseID,
		"is_active": true,
	}, nil)
	if err != nil {
		return nil, err
	}

	platformClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, refreshTokenIsPresent := cookies[platformClient["realm_name"]].(string)
	if refreshTokenIsPresent {
		refreshToken, err := jwt.DecodeJWTUnverified(refreshTokenStr)
		if err != nil {
			return nil, err
		}

		unexpired := jwt.VerifyExpired(refreshToken.Claims)
		if unexpired {
			kcUserID := jwt.GetKeycloakUserID(refreshToken.Claims)
			if kcUserID == "" {
				return nil, errors.New("invalid refresh token")

			}

			platformTokenResp, err := u.deps.KcTokenService.RefreshToken(
				ctx,
				platformClient["realm_name"],
				platformClient["client_id"],
				refreshTokenStr,
			)
			if err != nil {
				return nil, err
			}

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

			return &EnterpriseAuthToken{
				Token:             platformTokenResp,
				EnterpriseAccount: nil,
				RealmName:         platformClient["realm_name"],
				Permissions:       permissions,
			}, nil
		}
	}

	refreshTokenStr, ok := cookies[enterprise.KcRealmName].(string)
	if !ok {
		return nil, wraperror.NewValidationError(map[string]interface{}{
			"cookies[" + enterprise.KcRealmName + "]": utils.ErrorInputRequired,
		}, nil)
	}

	refreshToken, err := jwt.DecodeJWTUnverified(refreshTokenStr)
	if err != nil {
		return nil, err
	}

	kcUserID := jwt.GetKeycloakUserID(refreshToken.Claims)
	if kcUserID == "" {
		return nil, errors.New("invalid refresh token")
	}

	enterpriseTokenResp, err := u.deps.KcTokenService.RefreshToken(
		ctx,
		enterprise.KcRealmName,
		enterprise.KcClientID,
		refreshTokenStr,
	)
	if err != nil {
		return nil, err
	}

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

	return &EnterpriseAuthToken{
		Token:             enterpriseTokenResp,
		EnterpriseAccount: &enterpriseAccount,
		RealmName:         enterprise.KcRealmName,
		Permissions:       permissions,
	}, nil
}
