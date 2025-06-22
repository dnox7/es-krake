package usecase

import (
	"context"
	"strings"

	"github.com/dpe27/es-krake/internal/infrastructure/keycloak/utils"
	"github.com/dpe27/es-krake/pkg/jwt"
	baseUtils "github.com/dpe27/es-krake/pkg/utils"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

func (u *AuthUsecase) RegisterSns(ctx context.Context, kcUserInput map[string]interface{}) (*BuyerAuthToken, error) {
	buyerClient := u.deps.KcClientService.GetPlatformClient()
	realmName := buyerClient["realm_name"]
	clientID := buyerClient["client_id"]

	masterRealmToken, err := u.deps.KcTokenService.GetMasterRealmToken(ctx)
	if err != nil {
		return nil, err
	}

	kcUserTokenResp, err := u.deps.KcTokenService.GetTokenWithCode(
		ctx,
		realmName,
		clientID,
		kcUserInput["code"].(string),
		kcUserInput["redirect_uri"].(string),
	)
	if err != nil {
		if strings.Contains(err.Error(), utils.KeycloakInvalidAuthorizationCode) {
			return nil, wraperror.NewValidationError(
				map[string]interface{}{
					"code": baseUtils.ErrorInputFail,
				},
				err,
			)
		}

		if strings.Contains(err.Error(), utils.KeycloakInvalidRedirectURI) {
			return nil, wraperror.NewValidationError(
				map[string]interface{}{
					"redirect_uri": baseUtils.ErrorInputFail,
				}, err,
			)
		}
		return nil, err
	}

	kcUserToken, err := jwt.DecodeJWTUnverified(kcUserTokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	kcUserID := jwt.GetKeycloakUserID(kcUserToken.Claims)

	
}
