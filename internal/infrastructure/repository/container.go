package repository

import (
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	productRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	ProductContainer productRepo.RepositoryContainer
}

func NewRepositoriesContainer(pg *rdb.PostgreSQL) *RepositoriesContainer {
	return &RepositoriesContainer{
		ProductContainer: productRepo.NewRepositoryContainer(pg),
	}
}
