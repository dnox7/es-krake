package query

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func GetRoleType(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["role_type"],
		Name: "get_role_type",
		Args: graphql.FieldConfigArgument{
			"role_type_id": &graphql.ArgumentConfig{
				Type: output.AnyInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return authUsecase.GetRoleTypeByID(
				params.Context,
				params.Args["role_type_id"].(int),
			)
		},
	}
}
