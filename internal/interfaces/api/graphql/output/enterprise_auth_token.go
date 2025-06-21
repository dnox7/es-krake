package output

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func EnterpriseAuthTokenOutput(
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "enterprise_auth_token",
		Fields: graphql.Fields{
			"access_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).Token.AccessToken, nil
				},
			},
			"refresh_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).Token.RefreshToken, nil
				},
			},
			"refresh_expires_in": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).Token.RefreshExpiresIn, nil
				},
			},
			"realm_name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).RealmName, nil
				},
			},
			"permissions": &graphql.Field{
				Type: graphql.NewList(outputTypes["permission"]),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).Permissions, nil
				},
			},
			"enterprise_account": &graphql.Field{
				Type: outputTypes["enterprise_account"],
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.EnterpriseAuthToken).EnterpriseAccount, nil
				},
			},
		},
	})
}
