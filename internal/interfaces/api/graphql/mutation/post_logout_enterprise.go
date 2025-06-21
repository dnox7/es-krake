package mutation

import (
	"encoding/json"

	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostLogoutEnterprise(
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: output.Void,
		Name: "post_logout_enterprise",
		Args: graphql.FieldConfigArgument{
			"enterprise_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"kc_user_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"cookies": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			enterpriseID := params.Source.(map[string]interface{})["enterprise_id"].(string)
			kcUserID := params.Source.(map[string]interface{})["kc_user_id"].(string)
			cookies := params.Source.(map[string]interface{})["cookies"].(string)
			var cookiesMap map[string]interface{}
			err := json.Unmarshal([]byte(cookies), &cookiesMap)
			if err != nil {
				return nil, err
			}

			return nil, authUsecase.LogoutEnterprise(params.Context, enterpriseID, kcUserID, cookiesMap)
		},
	}
}
