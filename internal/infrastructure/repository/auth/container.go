package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	ActionRepository     domainRepo.ActionRepository
	PermissionRepository domainRepo.PermissionRepository
	ResourceRepository   domainRepo.ResourceRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		ActionRepository:     NewActionRepository(pg),
		PermissionRepository: NewPermissionRepository(pg),
		ResourceRepository:   NewResourceRepository(pg),
	}
}
