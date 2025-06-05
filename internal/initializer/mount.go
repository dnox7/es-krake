package initializer

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	"github.com/dpe27/es-krake/internal/interfaces/graphql"
	"github.com/dpe27/es-krake/internal/interfaces/handler"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/gin-gonic/gin"
)

func MountAll(pg *rdb.PostgreSQL, router *gin.Engine) error {
	repositories := repository.NewRepositoriesContainer(pg)
	services := service.NewServicesContainer(repositories)
	usecases := usecase.NewUsecasesContainer(repositories, services)

	schema, err := graphql.NewGraphQLSchema(&usecases)
	if err != nil {
		return fmt.Errorf("failed to create GraphQL schema: %v", err)
	}

	httpHandler := handler.NewHTTPHandler(schema)

	routerPF := router.Group("/pf")

	return nil
}
