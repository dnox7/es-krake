package repository

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/cart/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
)

type RepositoryContainer struct {
	CartRepo       domainRepo.CartRepository
	CartItemRepo   domainRepo.CartItemRepository
	CartStatusRepo domainRepo.CartStatusRepository
}

func NewRepositoryContainer(pg *rdb.PostgreSQL) RepositoryContainer {
	return RepositoryContainer{
		CartRepo:       NewCartRepository(pg),
		CartItemRepo:   NewCartItemRepository(pg),
		CartStatusRepo: NewCartStatusRepository(pg),
	}
}
