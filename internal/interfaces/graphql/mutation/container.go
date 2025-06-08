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
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})
}
