package graphql

import (
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/interfaces/graphql/mutation"
	"github.com/dpe27/es-krake/internal/interfaces/graphql/output"
	"github.com/dpe27/es-krake/internal/interfaces/graphql/query"

	"github.com/graphql-go/graphql"
)

func NewGraphQLSchema(
	repositories *repository.RepositoriesContainer,
) (graphql.Schema, error) {
	outputTypes := output.NewOutputTypes(repositories)

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query.NewQueriesContainer(repositories, outputTypes),
		Mutation: mutation.NewMutationsContainer(repositories, outputTypes),
	})
}
