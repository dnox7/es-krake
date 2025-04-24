package initializer

import (
	"pech/es-krake/internal/infrastructure/rdb"
	"pech/es-krake/internal/infrastructure/repository"
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
