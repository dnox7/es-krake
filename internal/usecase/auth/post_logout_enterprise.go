package usecase

import (
	"context"
	"net/http"

	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) PostLogoutEnterprise(
	ctx context.Context,
	enterpriseID string,
	kcUserID string,
	cookies map[string]interface{},
) error {
	enterprise, err := u.deps.EnterpriseRepo.TakeByConditions(ctx, map[string]interface{}{
		"id":        enterpriseID,
		"is_active": true,
	}, nil)
	if err != nil {
		return err
	}

	platformClient := u.deps.KcClientService.GetPlatformClient()
	refreshTokenStr, ok := cookies[platformClient["realm_name"]].(string)
	if ok {
		platformAccount, err := u.deps.PlatformAccountRepo.TakeByConditions(ctx, map[string]interface{}{
			"kc_user_id": kcUserID,
		}, nil)
		if err != nil {
			return err
		}

		_, err = u.deps.PlatformAccountEnterpriseAccessRepo.TakeByConditions(ctx, map[string]interface{}{
			"platform_account_id": platformAccount.ID,
			"enabled":             true,
			"enterprise_id":       enterprise.ID,
		}, nil)
		if err != nil {
			return err
		}

		return u.deps.KcUserService.LogoutOIDC(
			ctx,
			platformClient["realm_name"],
			platformClient["client_id"],
			refreshTokenStr,
		)
	}

	refreshTokenStr, ok = cookies[enterprise.KcRealmName].(string)
	if !ok {
		return wraperror.NewAPIError(
			http.StatusBadRequest,
			map[string]interface{}{
				"cookies[" + enterprise.KcRealmName + "]": utils.ErrorInputFail,
			}, nil,
		)
	}

	return u.deps.KcUserService.LogoutOIDC(
		ctx,
		enterprise.KcRealmName,
		enterprise.KcClientID,
		refreshTokenStr,
	)
}
