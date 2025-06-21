package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) LogoutBuyer(ctx context.Context, cookies map[string]interface{}) error {
	buyerClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, ok := cookies[buyerClient["realm_name"]].(string)
	if !ok {
		return wraperror.NewAPIError(
			http.StatusBadRequest,
			map[string]interface{}{
				"cookies[" + buyerClient["realm_name"] + "]": utils.ErrorInputFail,
			},
			errors.New("cookie without refresh token"),
		)
	}

	return u.deps.KcUserService.LogoutOIDC(
		ctx,
		buyerClient["realm_name"],
		buyerClient["client_id"],
		refreshTokenStr,
	)
}
