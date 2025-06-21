package mutation

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostLoginBuyer(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["buyer_auth_token"],
		Name: "post_login_buyer",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email := params.Source.(map[string]interface{})["email"].(string)
			password := params.Source.(map[string]interface{})["password"].(string)

			return authUsecase.LoginBuyer(params.Context, email, password)
		},
	}
}
