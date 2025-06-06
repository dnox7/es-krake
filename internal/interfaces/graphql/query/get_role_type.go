package query

import (
	"github.com/dpe27/es-krake/internal/interfaces/graphql/output"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/dpe27/es-krake/pkg/log"
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
			log.Debug(params.Context, "get role_type")
			return authUsecase.GetRoleTypeByID(
				params.Context,
				params.Args["role_type_id"].(int),
			)
		},
	}
}
