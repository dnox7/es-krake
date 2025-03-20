package repository

import (
	"pech/es-krake/internal/infrastructure/db"
	productRepo "pech/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	ProductContainer productRepo.RepositoryContainer
}

func NewRepositoriesContainer(pg *db.PostgreSQL) *RepositoriesContainer {
	return &RepositoriesContainer{
		ProductContainer: productRepo.NewRepositoryContainer(pg),
	}
}
