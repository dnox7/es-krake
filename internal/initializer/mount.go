package initializer

import (
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/internal/infrastructure/repository"
)

func MountAll(
	repositories *repository.RepositoriesContainer,
	pg *db.PostgreSQL,
) error {

	// graphqlSchema, err := graphql.NewGraphQLSchema(repositories)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create GraphQL schema: %v", err)
	// }
	return nil
}
