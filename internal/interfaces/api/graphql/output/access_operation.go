package output

import (
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func AccessOperationOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "access_operation",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.AccessOperation).ID, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.AccessOperation).Name, nil
				},
			},
			"code": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.AccessOperation).Code, nil
				},
			},
			"description": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.AccessOperation).Description, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					createdAt := params.Source.(entity.AccessOperation).CreatedAt
					return utils.ToDateTimeSQL(&createdAt), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updatedAt := params.Source.(entity.AccessOperation).UpdatedAt
					return utils.ToDateTimeSQL(&updatedAt), nil
				},
			},
		},
	})
}
