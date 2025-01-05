package container

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewGraphQLSchema(
	repositories *RepositoryContainers,
	services *ServiceContainers,
	db *gorm.DB,
	logger *logrus.Logger,
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
