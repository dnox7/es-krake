package output

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func BuyerAuthTokenOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "buyer_auth_token",
		Fields: graphql.Fields{
			"access_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.BuyerAuthToken).Token.AccessToken, nil
				},
			},
			"refresh_token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.BuyerAuthToken).Token.RefreshToken, nil
				},
			},
			"refresh_expires_in": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.BuyerAuthToken).Token.RefreshExpiresIn, nil
				},
			},
			"realm_name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(*usecase.BuyerAuthToken).RealmName, nil
				},
			},
		},
	})
}
