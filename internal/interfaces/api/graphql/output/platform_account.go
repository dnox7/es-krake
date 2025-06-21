package output

import (
	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	authUC "github.com/dpe27/es-krake/internal/usecase/auth"
	platformUC "github.com/dpe27/es-krake/internal/usecase/platform"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/graphql-go/graphql"
)

func PlatformAccountOutput(
	types map[string]*graphql.Object,
	authUsecase authUC.AuthUsecase,
	platformUsecase platformUC.PlatformUsecase,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "platform_account",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).ID, nil
				},
			},
			"kc_user_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).KcUserID, nil
				},
			},
			"role_id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).RoleID, nil
				},
			},
			"department_id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).DepartmentID, nil
				},
			},
			"has_password": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).HasPassword, nil
				},
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.PlatformAccount).Notes, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					created_at := params.Source.(entity.PlatformAccount).CreatedAt
					return utils.ToDateTimeSQL(&created_at), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updated_at := params.Source.(entity.PlatformAccount).UpdatedAt
					return utils.ToDateTimeSQL(&updated_at), nil
				},
			},
			"role": &graphql.Field{
				Type: types["role"],
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return authUsecase.GetRoleByIDWithPermissions(
						params.Context,
						params.Source.(entity.PlatformAccount).RoleID,
					)
				},
			},
			"department": &graphql.Field{
				Type: types["department"],
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return platformUsecase.GetDepartmentByID(
						params.Context,
						params.Source.(entity.PlatformAccount).DepartmentID,
					)
				},
			},
		},
	})
}
