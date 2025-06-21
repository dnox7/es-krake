package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) LogoutPlatform(ctx context.Context, cookies map[string]interface{}) error {
	platformClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, ok := cookies[platformClient["realm_name"]].(string)
	if !ok {
		return wraperror.NewAPIError(
			http.StatusBadRequest,
			map[string]interface{}{
				"cookies[" + platformClient["realm_name"] + "]": utils.ErrorInputFail,
			},
			errors.New("cookie without refresh token"),
		)
	}

	return u.deps.KcUserService.LogoutOIDC(
		ctx,
		platformClient["realm_name"],
		platformClient["client_id"],
		refreshTokenStr,
	)
}
