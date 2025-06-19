package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/seller/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	SellerAccountRepo domainRepo.SellerAccountRepository
	SellerInfoRepo    domainRepo.SellerInfoRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		SellerAccountRepo: NewSellerAccountRepository(pg),
		SellerInfoRepo:    NewSellerInfoRepository(pg),
	}
}
