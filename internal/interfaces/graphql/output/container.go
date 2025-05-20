package output

import (
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/graphql-go/graphql"
)

func NewOutputTypes(
	repositories *repository.RepositoriesContainer,
) map[string]*graphql.Object {
	outputTypes := make(map[string]*graphql.Object)

	for _, graphqlType := range []*graphql.Object{} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return outputTypes
}
