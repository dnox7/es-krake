package usecase

import (
	"context"
	"strings"

	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak/utils"
	"github.com/dpe27/es-krake/pkg/jwt"
	"github.com/dpe27/es-krake/pkg/wraperror"
)

const ErrorBuyerLoginFailed = "failed to login buyer"

type BuyerAuthToken struct {
	Token     kcdto.TokenEndpointResp
	RealmName string
}

func (u *AuthUsecase) LoginBuyer(ctx context.Context, username, password string) (*BuyerAuthToken, error) {
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
					"credentials": ErrorBuyerLoginFailed,
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
	_, err = u.deps.BuyerAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"kc_user_id":    kcUserID,
		"login_enabled": true,
		"mail_verified": true,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &BuyerAuthToken{
		Token:     tokenResp,
		RealmName: realmName,
	}, nil
}
