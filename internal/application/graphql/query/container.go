package query

import (
	"github.com/dpe27/es-krake/internal/infrastructure/repository"

	"github.com/graphql-go/graphql"
)

func NewQueriesContainer(
	repositories *repository.RepositoriesContainer,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})
}
