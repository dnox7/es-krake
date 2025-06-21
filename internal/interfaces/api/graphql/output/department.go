package output

import (
	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func DepartmentOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "department",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Department).ID, nil
				},
			},
			"code": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Department).Code, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Department).Name, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					created_at := params.Source.(entity.Department).CreatedAt
					return utils.ToDateTimeSQL(&created_at), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updated_at := params.Source.(entity.Department).UpdatedAt
					return utils.ToDateTimeSQL(&updated_at), nil
				},
			},
		},
	})
}
