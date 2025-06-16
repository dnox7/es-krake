package output

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func PlatformLogin(
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "platform_login",
		Fields: graphql.Fields{
			"access_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.LoginPlatformResponse).Token.AccessToken, nil
				},
			},
			"refresh_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.LoginPlatformResponse).Token.RefreshToken, nil
				},
			},
			"refresh_expires_in": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.LoginPlatformResponse).Token.RefreshExpiresIn, nil
				},
			},
			"permissions": &graphql.Field{
				Type: graphql.NewList(outputTypes["permission"]),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.LoginPlatformResponse).Permissions, nil
				},
			},
		},
	})
}
