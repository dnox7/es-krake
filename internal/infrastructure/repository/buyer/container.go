package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/buyer/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	BuyerAccountRepo domainRepo.BuyerAccountRepository
	BuyerInfoRepo    domainRepo.BuyerInfoRepository
	GenderRepo       domainRepo.GenderRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		BuyerAccountRepo: NewBuyerAccountRepository(pg),
		BuyerInfoRepo:    NewBuyerInfoRepository(pg),
		GenderRepo:       NewGenderRepository(pg),
	}
}
