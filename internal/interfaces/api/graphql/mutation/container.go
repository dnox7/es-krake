package mutation

import (
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/graphql-go/graphql"
)

func NewMutationsContainer(
	usecases *usecase.UsecasesContainer,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"post_login_platform": PostLoginPlatform(
				outputTypes,
				usecases.AuthUsecase,
			),
			"post_login_enterprise": PostLoginEnterprise(
				outputTypes,
				usecases.AuthUsecase,
			),
		},
	})
}
