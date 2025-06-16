package usecase

import (
	authRepositories "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/auth/service"
	enterpriseRepositories "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	platformRepositories "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
)

type AuthUsecaseDeps struct {
	RoleTypeRepo   authRepositories.RoleTypeRepository
	PermissionRepo authRepositories.PermissionRepository
	RoleRepo       authRepositories.RoleRepository

	PlatformAccountRepo                 platformRepositories.PlatformAccountRepository
	PlatformAccountEnterpriseAccessRepo platformRepositories.PlatformAccountEnterpriseAccessRepository

	EnterpriseRepo        enterpriseRepositories.EnterpriseRepository
	EnterpriseAccountRepo enterpriseRepositories.EnterpriseAccountRepository

	PermissionService      service.PermissionService
	AccessOperationService service.AccessOperationService

	KcTokenService  keycloak.KcTokenService
	KcClientService keycloak.KcClientService

	Cache redis.RedisRepository
}

type AuthUsecase struct {
	logger *log.Logger
	deps   *AuthUsecaseDeps
}

func NewAuthUsecase(deps *AuthUsecaseDeps) AuthUsecase {
	return AuthUsecase{
		logger: log.With("object", "auth_usecase"),
		deps:   deps,
	}
}
