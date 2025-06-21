package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) RefreshTokenBuyer(ctx context.Context, cookies map[string]interface{}) (*BuyerAuthToken, error) {
	buyerClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, ok := cookies[buyerClient["realm_name"]].(string)
	if !ok {
		return nil, wraperror.NewAPIError(
			http.StatusBadRequest,
			map[string]interface{}{
				"cookies[" + buyerClient["realm_name"] + "]": utils.ErrorInputFail,
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
		buyerClient["realm_name"],
		buyerClient["client_id"],
		refreshTokenStr,
	)
	if err != nil {
		return nil, err
	}

	_, err = u.deps.BuyerAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"kc_user_id":     kcUserID,
		"login_enable":   true,
		"email_verified": true,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &BuyerAuthToken{
		Token:     tokenResp,
		RealmName: buyerClient["realm_name"],
	}, nil

}
