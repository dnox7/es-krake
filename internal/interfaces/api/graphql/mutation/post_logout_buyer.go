package mutation

import (
	"encoding/json"

	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostLogoutBuyer(
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: output.Void,
		Name: "post_logout_buyer",
		Args: graphql.FieldConfigArgument{
			"cookies": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			cookies := params.Source.(map[string]interface{})["cookies"].(string)
			var cookiesMap map[string]interface{}
			err := json.Unmarshal([]byte(cookies), &cookiesMap)
			if err != nil {
				return nil, err
			}

			return nil, authUsecase.LogoutBuyer(params.Context, cookiesMap)
		},
	}
}
