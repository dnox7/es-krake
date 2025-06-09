package query

import (
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/graphql-go/graphql"
)

func NewQueriesContainer(
	usecases *usecase.UsecasesContainer,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"get_role_type": GetRoleType(
				outputTypes,
				usecases.AuthUsecase,
			),
		},
	})
}
