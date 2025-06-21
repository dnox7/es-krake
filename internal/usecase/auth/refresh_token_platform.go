package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) RefreshTokenPlatform(ctx context.Context, cookies map[string]interface{}) (*PlatformAuthToken, error) {
	platformClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, ok := cookies[platformClient["realm_name"]].(string)
	if !ok {
		return nil, wraperror.NewAPIError(
			http.StatusBadRequest,
			map[string]interface{}{
				"cookies": utils.ErrorInputFail,
			},
			errors.New("cookie without refresh token"),
		)
	}

	refreshToken, err := jwt.DecodeJWTUnverified(refreshTokenStr)
	if err != nil {
		return nil, wraperror.NewAPIError(
			http.StatusBadRequest,
			err.Error(),
			err,
		)
	}

	kcUserID := jwt.GetKeycloakUserID(refreshToken.Claims)
	tokenResp, err := u.deps.KcTokenService.RefreshToken(
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

	permissions, err := u.deps.PermissionService.GetPermissionsWithRoleID(ctx, platformAccount.RoleID)
	if err != nil {
		return nil, err
	}

	return &PlatformAuthToken{
		Token:       tokenResp,
		RealmName:   platformClient["realm_name"],
		Permissions: permissions,
	}, nil
}
