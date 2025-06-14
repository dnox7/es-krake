package output

import (
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func RoleTypeOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "role_type",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: AnyInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.RoleType).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.RoleType).Name, nil
					},
				},
				"created_at": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						created_at := params.Source.(entity.RoleType).CreatedAt
						return utils.ToDateTimeSQL(&created_at), nil
					},
				},
				"updated_at": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						updated_at := params.Source.(entity.RoleType).UpdatedAt
						return utils.ToDateTimeSQL(&updated_at), nil
					},
				},
			}
		}),
	})
}
