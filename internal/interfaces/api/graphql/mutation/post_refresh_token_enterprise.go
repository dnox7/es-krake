package mutation

import (
	"encoding/json"

	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostRefreshTokenEnterprise(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["enterprise_auth_token"],
		Name: "post_refresh_token_enterprise",
		Args: graphql.FieldConfigArgument{
			"enterprise_id": &graphql.ArgumentConfig{
				Type: output.AnyInt,
			},
			"cookies": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			enterpriseID := params.Source.(map[string]interface{})["enterprise_id"].(int)
			jsonCookies := params.Source.(map[string]interface{})["cookies"].(string)
			var cookies map[string]interface{}
			err := json.Unmarshal([]byte(jsonCookies), &cookies)
			if err != nil {
				return nil, err
			}

			return authUsecase.RefreshTokenEnterprise(params.Context, enterpriseID, cookies)
		},
	}
}
