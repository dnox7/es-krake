package usecase

import (
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	authUC "github.com/dpe27/es-krake/internal/usecase/auth"
	batchLogUC "github.com/dpe27/es-krake/internal/usecase/batchlog"
)

type UsecasesContainer struct {
	AuthUsecase     authUC.AuthUsecase
	BatchLogUsecase batchLogUC.BatchLogUsecase
}

func NewUsecasesContainer(
	repositories *repository.RepositoriesContainer,
	services *service.ServicesContainer,
	redisRepo redis.RedisRepository,
) UsecasesContainer {
	return UsecasesContainer{
		AuthUsecase: authUC.NewAuthUsecase(&authUC.AuthUsecaseDeps{
			RoleTypeRepo: repositories.AuthContainer.RoleTypeRepo,
			Cache:        redisRepo,
		}),
		BatchLogUsecase: batchLogUC.NewBatchLogUsecase(&batchLogUC.BatchLogUsecaseDeps{
			BatchLogRepo:     repositories.BatchLogContainer.BatchLogRepo,
			BatchLogTypeRepo: repositories.BatchLogContainer.BatchLogTypeRepo,
			Cache:            redisRepo,
		}),
	}
}
