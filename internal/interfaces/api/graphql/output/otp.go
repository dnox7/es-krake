package output

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func OtpStatusOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "otp_status",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(usecase.OtpStatus).Status, nil
				},
			},
			"email_verified": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(usecase.OtpStatus).EmailVerified, nil
				},
			},
			"is_registered_password": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(usecase.OtpStatus).IsRegisteredPassword, nil
				},
			},
			"is_registered_sns": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(usecase.OtpStatus).IsRegisteredSns, nil
				},
			},
		},
	})
}
