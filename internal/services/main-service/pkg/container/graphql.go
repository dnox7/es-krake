package container

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewGraphQLSchema(
	repositories *RepositoryContainers,
	services *ServiceContainers,
	db *gorm.DB,
) (graphql.Schema, error) {
	outputTypes := make(map[string]*graphql.Object)
	for _, graphqlType := range []*graphql.Object{} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields{},
		}),
	})
}
