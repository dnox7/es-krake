package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	DepartmentRepo                      domainRepo.DepartmentRepository
	PlatformAccountRepo                 domainRepo.PlatformAccountRepository
	PlatformAccountEnterpriseAccessRepo domainRepo.PlatformAccountEnterpriseAccessRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		DepartmentRepo:                      NewDepartmentRepository(pg),
		PlatformAccountRepo:                 NewPlatformAccountRepository(pg),
		PlatformAccountEnterpriseAccessRepo: NewPlatformAccountEnterpriseAccessRepo(pg),
	}
}
