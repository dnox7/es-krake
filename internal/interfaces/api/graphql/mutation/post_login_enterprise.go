package mutation

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostLoginEnterprise(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["enterprise_auth_token"],
		Name: "post_login_enterprise",
		Args: graphql.FieldConfigArgument{
			"enterprise_id": &graphql.ArgumentConfig{
				Type: output.AnyInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email := params.Source.(map[string]interface{})["email"].(string)
			password := params.Source.(map[string]interface{})["password"].(string)
			enterpriseID := params.Args["enterprise_id"].(int)

			return authUsecase.LoginEnterprise(params.Context, enterpriseID, email, password)
		},
	}
}
