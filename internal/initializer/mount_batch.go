package initializer

import (
	"github.com/dpe27/es-krake/config"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/notify"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	"github.com/dpe27/es-krake/internal/interfaces/batch/handler"
	"github.com/dpe27/es-krake/internal/interfaces/batch/jobs"
	"github.com/dpe27/es-krake/internal/interfaces/batch/router"
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/gin-gonic/gin"
)

func MountBatch(
	cfg *config.Config,
	pg *rdb.PostgreSQL,
	mongo *mdb.Mongo,
	redisRepo redis.RedisRepository,
	notifier notify.DiscordNotifier,
	ginEngine *gin.Engine,
) error {
	repositories := repository.NewRepositoriesContainer(pg, mongo)
	services := service.NewServicesContainer(repositories, redisRepo)
	usecases := usecase.NewUsecasesContainer(repositories, services, redisRepo)
	batchContainer := jobs.NewBatchContainer(&usecases)

	batchHandler := handler.NewBatchHandler(
		batchContainer, notifier,
		cfg.App.LogLevel == "DEBUG",
	)

	router.BindBatchRoutes(ginEngine.Group("/"), batchHandler)
	return nil
}
