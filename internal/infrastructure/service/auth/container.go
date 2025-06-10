package service

import (
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	repository "github.com/dpe27/es-krake/internal/infrastructure/repository/auth"
)

type ServiceContainer struct {
	AccessOperationService domainService.AccessOperationService
	PermissionService      domainService.PermissionService
}

func NewServiceContainer(
	repos repository.RepositoryContainer,
	cache redis.RedisRepository,
) ServiceContainer {
	return ServiceContainer{
		AccessOperationService: NewAccessOperationService(
			repos.AccessOperationRepo,
			cache,
		),
		PermissionService: NewPermissionService(
			repos.PermissionRepo,
		),
	}
}
