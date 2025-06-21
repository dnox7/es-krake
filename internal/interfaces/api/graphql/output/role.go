package output

import (
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func RoleOutput(
	types map[string]*graphql.Object,
	usecase usecase.AuthUsecase,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "role",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Role).ID, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Role).Name, nil
				},
			},
			"role_type_id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Role).RoleTypeID, nil
				},
			},
			"display_order": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Role).DisplayOrder, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					created_at := params.Source.(entity.Role).CreatedAt
					return utils.ToDateTimeSQL(&created_at), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updated_at := params.Source.(entity.Role).UpdatedAt
					return utils.ToDateTimeSQL(&updated_at), nil
				},
			},
			"role_type": &graphql.Field{
				Type: types["role_type"],
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return usecase.GetRoleTypeByID(params.Context, params.Source.(entity.Role).RoleTypeID)
				},
			},
			"permissions": &graphql.Field{
				Type: graphql.NewList(types["permissions"]),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return usecase.GetPermissionsWithRoleID(params.Context, params.Source.(entity.Role).ID)
				},
			},
		},
	})
}
