package query

import (
	usecase "github.com/dpe27/es-krake/internal/usecase/auth"
	"github.com/graphql-go/graphql"
)

func GetBuyerOtpStatus(
	outputTypes map[string]*graphql.Object,
	authUsecase usecase.AuthUsecase,
) *graphql.Field {
	return &graphql.Field{
		Type: outputTypes["otp_status"],
		Name: "get_buyer_otp_status",
		Args: graphql.FieldConfigArgument{
			"otp": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			otp := params.Args["otp"].(string)
			return authUsecase.CheckBuyerOtpStatus(
				params.Context,
				otp,
			)
		},
	}
}
