package mutation

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PostSignupBuyer(
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: output.Void,
		Name: "post_signup_buyer",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			kcUserInput := params.Source.(map[string]interface{})["keycloak_user"].(map[string]interface{})

			return nil, authUsecase.SignupBuyer(params.Context, kcUserInput)
		},
	}
}
