package initializer

import (
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
)

func MountAll(
	repositories *repository.RepositoriesContainer,
	pg *rdb.PostgreSQL,
) error {

	// graphqlSchema, err := graphql.NewGraphQLSchema(repositories)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create GraphQL schema: %v", err)
	// }
	return nil
}
