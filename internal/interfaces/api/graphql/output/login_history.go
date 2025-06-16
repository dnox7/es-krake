package output

import (
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/graphql-go/graphql"
)

func LoginHistoryOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "login_history",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).ID, nil
				},
			},
			"kc_user_id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).KcUserID, nil
				},
			},
			"logged_in_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).LoggedInAt, nil
				},
			},
			"logged_device": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).LoggedDevice, nil
				},
			},
			"logged_ip_address": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).LoggedIPAddress, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).CreatedAt, nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.LoginHistory).UpdatedAt, nil
				},
			},
		},
	})
}
