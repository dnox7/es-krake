package service

import (
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	repository "github.com/dpe27/es-krake/internal/infrastructure/repository/auth"
)

type ServiceContainer struct {
	AccessOperationService domainService.AccessOperationService
	PermissionService      domainService.PermissionService
}

func NewServiceContainer(repos repository.RepositoryContainer) ServiceContainer {
	return ServiceContainer{
		AccessOperationService: NewAccessOperationService(
			repos.AccessRequirementRepo,
		),
		PermissionService: NewPermissionService(
			repos.PermissionRepo,
		),
	}
}
