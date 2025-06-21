package initializer

import (
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/aws"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	"github.com/dpe27/es-krake/internal/interfaces/api/graphql"
	"github.com/dpe27/es-krake/internal/interfaces/api/handler"
	"github.com/dpe27/es-krake/internal/interfaces/api/router"
	"github.com/dpe27/es-krake/internal/interfaces/middleware"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/gin-gonic/gin"
)

func MountAPI(
	cfg *config.Config,
	pg *rdb.PostgreSQL,
	mongo *mdb.Mongo,
	redisRepo redis.RedisRepository,
	ses aws.SesService,
	ginEngine *gin.Engine,
) error {
	inputValidator, err := validator.NewJsonSchemaValidator(cfg)
	if err != nil {
		return fmt.Errorf("failed to create input validator: %v", err)
	}

	repositories := repository.NewRepositoriesContainer(pg, mongo)
	services := service.NewServicesContainer(cfg, ses, repositories, redisRepo)
	usecases := usecase.NewUsecasesContainer(pg, repositories, services, redisRepo)
	schema, err := graphql.NewGraphQLSchema(&usecases)
	if err != nil {
		return fmt.Errorf("failed to create GraphQL schema: %v", err)
	}

	httpHandler := handler.NewHTTPHandler(
		schema,
		cfg.App.LogLevel == "DEBUG",
		inputValidator,
	)

	authenMiddleware := middleware.NewAuthenMiddleware(cfg, services.KeycloakContainer.KeyService)
	permMiddleware := middleware.NewPermMiddleware(
		services.AuthContainer.AccessOperationService,
		repositories.AuthContainer.AccessRequirementRepo,
		repositories.EnterpriseContainer.EnterpriseAccountRepo,
		repositories.PlatformContainer.PlatformAccountRepo,
		services.AuthContainer.PermissionService,
		repositories.AuthContainer.RoleRepo,
	)

	router.BindPlatformRoute(ginEngine.Group("/pf"), httpHandler.Pf, authenMiddleware, permMiddleware)
	router.BindEnterpriseRoute(ginEngine.Group("/ent"), httpHandler.Ent, authenMiddleware, permMiddleware)
	router.BindBuyerRoute(ginEngine.Group("/buyer"), httpHandler.Buyer, authenMiddleware, permMiddleware)
	return nil
}
