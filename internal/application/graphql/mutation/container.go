package mutation

import (
	"github.com/dpe27/es-krake/internal/infrastructure/repository"

	"github.com/graphql-go/graphql"
)

func NewMutationsContainer(
	repositories *repository.RepositoriesContainer,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})
}
