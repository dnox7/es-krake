package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	AccessRequirementRepo domainRepo.AccessRequirementRepository
	ActionRepo            domainRepo.ActionRepository
	FunctionCodeRepo      domainRepo.FunctionCodeRepository
	PermissionRepo        domainRepo.PermissionRepository
	ResourceRepo          domainRepo.ResourceRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AccessRequirementRepo: NewAccessRequirementRepository(pg),
		ActionRepo:            NewActionRepository(pg),
		FunctionCodeRepo:      NewFunctionCodeRepository(pg),
		PermissionRepo:        NewPermissionRepository(pg),
		ResourceRepo:          NewResourceRepository(pg),
	}
}
