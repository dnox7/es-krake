package graphql

import (
	"pech/es-krake/internal/application/graphql/mutation"
	"pech/es-krake/internal/application/graphql/output"
	"pech/es-krake/internal/application/graphql/query"
	"pech/es-krake/internal/infrastructure/repository"

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
