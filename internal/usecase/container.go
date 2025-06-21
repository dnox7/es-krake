package usecase

import (
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
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
	pg *rdb.PostgreSQL,
	repositories *repository.RepositoriesContainer,
	services *service.ServicesContainer,
	redisRepo redis.RedisRepository,
) UsecasesContainer {
	return UsecasesContainer{
		AuthUsecase: authUC.NewAuthUsecase(&authUC.AuthUsecaseDeps{
			DB:                                  pg,
			RoleTypeRepo:                        repositories.AuthContainer.RoleTypeRepo,
			PermissionRepo:                      repositories.AuthContainer.PermissionRepo,
			RoleRepo:                            repositories.AuthContainer.RoleRepo,
			OtpRepo:                             repositories.AuthContainer.OtpRepo,
			PlatformAccountRepo:                 repositories.PlatformContainer.PlatformAccountRepo,
			BuyerAccountRepo:                    repositories.BuyerContainer.BuyerAccountRepo,
			EnterpriseAccountRepo:               repositories.EnterpriseContainer.EnterpriseAccountRepo,
			EnterpriseRepo:                      repositories.EnterpriseContainer.EnterpriseRepo,
			PlatformAccountEnterpriseAccessRepo: repositories.PlatformContainer.PlatformAccountEnterpriseAccessRepo,
			Cache:                               redisRepo,
			AccessOperationService:              services.AuthContainer.AccessOperationService,
			PermissionService:                   services.AuthContainer.PermissionService,
			KcTokenService:                      services.KeycloakContainer.TokenService,
			KcClientService:                     services.KeycloakContainer.ClientService,
			KcUserService:                       services.KeycloakContainer.UserService,
			MailService:                         services.MailService,
			BuyerService:                        services.BuyerContainer.BuyerService,
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
