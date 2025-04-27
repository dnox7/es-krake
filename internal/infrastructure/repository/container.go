package repository

import (
	"pech/es-krake/internal/infrastructure/rdb"
	productRepo "pech/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	ProductContainer productRepo.RepositoryContainer
}

func NewRepositoriesContainer(pg *rdb.PostgreSQL) *RepositoriesContainer {
	return &RepositoriesContainer{
		ProductContainer: productRepo.NewRepositoryContainer(pg),
	}
}
