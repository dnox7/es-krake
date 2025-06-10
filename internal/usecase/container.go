package usecase

import (
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	authUC "github.com/dpe27/es-krake/internal/usecase/auth"
)

type UsecasesContainer struct {
	AuthUsecase authUC.AuthUsecase
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
	}
}
