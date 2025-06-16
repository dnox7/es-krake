package mutation

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostLoginPlatform(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["platform_login"],
		Name: "post_login_platform",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email := params.Source.(map[string]interface{})["email"].(string)
			password := params.Source.(map[string]interface{})["password"].(string)

			return authUsecase.LoginPlatform(params.Context, email, password)
		},
	}
}
