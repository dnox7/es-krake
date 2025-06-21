package mutation

import (
	"encoding/json"

	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostRefreshTokenPlatform(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["platform_auth_token"],
		Name: "post_refresh_token_platform",
		Args: graphql.FieldConfigArgument{
			"cookies": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			jsonCookies := params.Source.(map[string]interface{})["cookies"].(string)
			var cookies map[string]interface{}
			err := json.Unmarshal([]byte(jsonCookies), &cookies)
			if err != nil {
				return nil, err
			}

			return authUsecase.RefreshTokenPlatform(params.Context, cookies)
		},
	}
}
