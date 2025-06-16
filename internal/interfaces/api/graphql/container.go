package graphql

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/mutation"
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/output"
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql/query"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/graphql-go/graphql"
)

func NewGraphQLSchema(usecases *usecase.UsecasesContainer) (graphql.Schema, error) {
	outputTypes := output.NewOutputTypes(usecases)

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query.NewQueriesContainer(usecases, outputTypes),
		Mutation: mutation.NewMutationsContainer(usecases, outputTypes),
	})
}
