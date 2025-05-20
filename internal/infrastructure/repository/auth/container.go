package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	AccessOperationRepo            domainRepo.AccessOperationRepository
	AccessRequirementRepo          domainRepo.AccessRequirementRepository
	AccessRequirementOperationRepo domainRepo.AccessRequirementOperationRepository
	ActionRepo                     domainRepo.ActionRepository
	PermissionRepo                 domainRepo.PermissionRepository
	PermissionOpeartionRepo        domainRepo.PermissionOperationRepository
	ResourceRepo                   domainRepo.ResourceRepository
	RoleRepo                       domainRepo.RoleRepository
	RolePermissionRepo             domainRepo.RolePermissionRepository
	RoleTypeRepo                   domainRepo.RoleTypeRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		AccessOperationRepo:            NewAccessOperationRepository(pg),
		AccessRequirementRepo:          NewAccessRequirementRepository(pg),
		AccessRequirementOperationRepo: NewAccessRequirementOperationRepository(pg),
		ActionRepo:                     NewActionRepository(pg),
		PermissionRepo:                 NewPermissionRepository(pg),
		PermissionOpeartionRepo:        NewPermissionOperationRepository(pg),
		ResourceRepo:                   NewResourceRepository(pg),
		RoleRepo:                       NewRoleRepository(pg),
		RolePermissionRepo:             NewRolePermissionRepository(pg),
		RoleTypeRepo:                   NewRoleTypeRepository(pg),
	}
}
