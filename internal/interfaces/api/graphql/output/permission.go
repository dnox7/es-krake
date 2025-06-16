package output

import (
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func PermissionOutput(
	types map[string]*graphql.Object,
	usecase usecase.AuthUsecase,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "permission",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Permission).ID, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Permission).Name, nil
				},
			},
			"access_operations": &graphql.Field{
				Type: graphql.NewList(types["access_operation"]),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					permID := params.Source.(entity.Permission).ID
					return usecase.GetOperationsByPermissionID(params.Context, permID)
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					createdAt := params.Source.(entity.Permission).CreatedAt
					return utils.ToDateTimeSQL(&createdAt), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updatedAt := params.Source.(entity.Permission).UpdatedAt
					return utils.ToDateTimeSQL(&updatedAt), nil
				},
			},
		},
	})
}
