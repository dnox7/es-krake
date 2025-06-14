package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/batchlog/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	BatchLogTypeRepo domainRepo.BatchLogTypeRepository
	BatchLogRepo     domainRepo.BatchLogRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		BatchLogTypeRepo: NewBatchLogTypeRepository(pg),
		BatchLogRepo:     NewBatchLogRepository(pg),
	}
}
