package usecase

import (
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/service"
	authUC "github.com/dpe27/es-krake/internal/usecase/auth"
	batchLogUC "github.com/dpe27/es-krake/internal/usecase/batchlog"
	platformUC "github.com/dpe27/es-krake/internal/usecase/platform"
)

type UsecasesContainer struct {
	AuthUsecase     authUC.AuthUsecase
	BatchLogUsecase batchLogUC.BatchLogUsecase
	PlatformUsecase platformUC.PlatformUsecase
}

func NewUsecasesContainer(
	repositories *repository.RepositoriesContainer,
	services *service.ServicesContainer,
	redisRepo redis.RedisRepository,
) UsecasesContainer {
	return UsecasesContainer{
		AuthUsecase: authUC.NewAuthUsecase(&authUC.AuthUsecaseDeps{
			RoleTypeRepo:           repositories.AuthContainer.RoleTypeRepo,
			PermissionRepo:         repositories.AuthContainer.PermissionRepo,
			RoleRepo:               repositories.AuthContainer.RoleRepo,
			PlatformAccountRepo:    repositories.PlatformContainer.PlatformAccountRepo,
			PermissionService:      services.AuthContainer.PermissionService,
			AccessOperationService: services.AuthContainer.AccessOperationService,
			KcTokenService:         services.KeycloakContainer.TokenService,
			KcClientService:        services.KeycloakContainer.ClientService,
			Cache:                  redisRepo,
		}),
		BatchLogUsecase: batchLogUC.NewBatchLogUsecase(&batchLogUC.BatchLogUsecaseDeps{
			BatchLogRepo:     repositories.BatchLogContainer.BatchLogRepo,
			BatchLogTypeRepo: repositories.BatchLogContainer.BatchLogTypeRepo,
			Cache:            redisRepo,
		}),
		PlatformUsecase: platformUC.NewPlatformUsecase(&platformUC.PlatformUsecaseDeps{
			PlatformAccountRepo: repositories.PlatformContainer.PlatformAccountRepo,
			DepartmentRepo:      repositories.PlatformContainer.DepartmentRepo,
			Cache:               redisRepo,
		}),
	}
}
