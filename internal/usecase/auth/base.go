package usecase

import (
	authRepositories "github.com/dpe27/es-krake/internal/domain/auth/repository"
	authServices "github.com/dpe27/es-krake/internal/domain/auth/service"
	buyerRepositories "github.com/dpe27/es-krake/internal/domain/buyer/repository"
	buyerServices "github.com/dpe27/es-krake/internal/domain/buyer/service"
	enterpriseRepositories "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	platformRepositories "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	mailService "github.com/dpe27/es-krake/internal/infrastructure/service/mail"
	"github.com/dpe27/es-krake/pkg/log"
)

type AuthUsecaseDeps struct {
	DB *rdb.PostgreSQL

	RoleTypeRepo   authRepositories.RoleTypeRepository
	PermissionRepo authRepositories.PermissionRepository
	RoleRepo       authRepositories.RoleRepository
	OtpRepo        authRepositories.OtpRepository

	PlatformAccountRepo                 platformRepositories.PlatformAccountRepository
	PlatformAccountEnterpriseAccessRepo platformRepositories.PlatformAccountEnterpriseAccessRepository

	EnterpriseRepo        enterpriseRepositories.EnterpriseRepository
	EnterpriseAccountRepo enterpriseRepositories.EnterpriseAccountRepository

	BuyerAccountRepo buyerRepositories.BuyerAccountRepository

	PermissionService      authServices.PermissionService
	AccessOperationService authServices.AccessOperationService
	BuyerService           buyerServices.BuyerService

	KcTokenService  keycloak.KcTokenService
	KcClientService keycloak.KcClientService
	KcUserService   keycloak.KcUserService

	MailService mailService.MailService
	Cache       redis.RedisRepository
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
