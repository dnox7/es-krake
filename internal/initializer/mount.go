package initializer

import (
	"fmt"

	"github.com/dpe27/es-krake/config"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	"github.com/dpe27/es-krake/internal/interfaces/graphql"
	"github.com/dpe27/es-krake/internal/interfaces/handler"
	"github.com/dpe27/es-krake/internal/interfaces/router"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/gin-gonic/gin"
)

func MountAll(
	cfg *config.Config,
	pg *rdb.PostgreSQL,
	mongo *mdb.Mongo,
	redisRepo redis.RedisRepository,
	ginEngine *gin.Engine,
) error {
	repositories := repository.NewRepositoriesContainer(pg, mongo)
	services := service.NewServicesContainer(repositories, redisRepo)
	usecases := usecase.NewUsecasesContainer(repositories, services, redisRepo)

	schema, err := graphql.NewGraphQLSchema(&usecases)
	if err != nil {
		return fmt.Errorf("failed to create GraphQL schema: %v", err)
	}

	debug := cfg.App.LogLevel == "DEBUG"
	httpHandler := handler.NewHTTPHandler(debug, schema)

	router.BindPlatformRoute(ginEngine.Group("/pf"), httpHandler.PF)
	return nil
}
