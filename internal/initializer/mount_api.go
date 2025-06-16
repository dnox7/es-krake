package initializer

import (
	"fmt"

	"github.com/dpe27/es-krake/config"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql"
	"github.com/dpe27/es-krake/internal/interfaces/api/handler"
	"github.com/dpe27/es-krake/internal/interfaces/api/router"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/gin-gonic/gin"
)

func MountAPI(
	cfg *config.Config,
	pg *rdb.PostgreSQL,
	mongo *mdb.Mongo,
	redisRepo redis.RedisRepository,
	ginEngine *gin.Engine,
) error {
	repositories := repository.NewRepositoriesContainer(pg, mongo)
	services := service.NewServicesContainer(cfg, repositories, redisRepo)
	usecases := usecase.NewUsecasesContainer(repositories, services, redisRepo)
	schema, err := graphql.NewGraphQLSchema(&usecases)
	if err != nil {
		return fmt.Errorf("failed to create GraphQL schema: %v", err)
	}

	httpHandler := handler.NewHTTPHandler(
		schema, cfg.App.LogLevel == "DEBUG",
	)

	router.BindPlatformRoute(ginEngine.Group("/pf"), httpHandler.PF)
	router.BindEnterpriseRoute(ginEngine.Group("/ent"), httpHandler.Ent)
	return nil
}
