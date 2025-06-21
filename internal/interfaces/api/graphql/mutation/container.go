package mutation

import (
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/graphql-go/graphql"
)

func NewMutationsContainer(
	usecases *usecase.UsecasesContainer,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"post_login_platform": PostLoginPlatform(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_login_enterprise": PostLoginEnterprise(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_login_buyer": PostLoginBuyer(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_signup_buyer": PostSignupBuyer(
				usecases.AuthUsecase,
			),
			"post_logout_buyer": PostLogoutBuyer(
				usecases.AuthUsecase,
			),
			"post_logout_platform": PostLogoutPlatform(
				usecases.AuthUsecase,
			),
			"post_logout_enterprise": PostLogoutEnterprise(
				usecases.AuthUsecase,
			),
			"post_refresh_token_buyer": PostRefreshTokenBuyer(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_refresh_token_platform": PostRefreshTokenPlatform(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_refresh_token_enterprise": PostRefreshTokenEnterprise(
				outputTypes,
				usecases.AuthUsecase,
			),
		},
	})
}
