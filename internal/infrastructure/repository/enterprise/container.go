package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	EnterpriseRepo        domainRepo.EnterpriseRepository
	EnterpriseAccountRepo domainRepo.EnterpriseAccountRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		EnterpriseRepo:        NewEnterpriseRepository(pg),
		EnterpriseAccountRepo: NewEnterpriseAccountRepository(pg),
	}
}
