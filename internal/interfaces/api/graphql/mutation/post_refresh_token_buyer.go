package mutation

import (
	"encoding/json"

	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostRefreshTokenBuyer(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["buyer_auth_token"],
		Name: "post_refresh_token_buyer",
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

			return authUsecase.RefreshTokenBuyer(params.Context, cookies)
		},
	}
}