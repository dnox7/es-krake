package output

import (
	"github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	authUC "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func EnterpriseAccountOutput(
	types map[string]*graphql.Object,
	authUsecase authUC.AuthUsecase,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "enterprise_account",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).ID, nil
				},
			},
			"kc_user_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).KcUserID, nil
				},
			},
			"role_id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).RoleID, nil
				},
			},
			"enterprise_id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).EnterpriseID, nil
				},
			},
			"has_password": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).HasPassword, nil
				},
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.EnterpriseAccount).Notes, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					createdAt := params.Source.(entity.EnterpriseAccount).CreatedAt
					return utils.ToDateTimeSQL(&createdAt), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updatedAt := params.Source.(entity.EnterpriseAccount).UpdatedAt
					return utils.ToDateTimeSQL(&updatedAt), nil
				},
			},
			"role": &graphql.Field{
				Type: types["role"],
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return authUsecase.GetRoleByIDWithPermissions(
						params.Context,
						params.Source.(entity.EnterpriseAccount).RoleID,
					)
				},
			},
		},
	})
}
