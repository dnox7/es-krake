package repository

import (
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	authRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/auth"
	enterpriseRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/enterprise"
	prodRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	AuthContainer       authRepo.RepositoryContainer
	EnterpriseContainer enterpriseRepo.RepositoryContainer
	ProductContainer    prodRepo.RepositoryContainer
}

func NewRepositoriesContainer(pg *rdb.PostgreSQL) *RepositoriesContainer {
	return &RepositoriesContainer{
		AuthContainer:       authRepo.NewRepositoryContainer(pg),
		EnterpriseContainer: enterpriseRepo.NewRepositoryContainer(pg),
		ProductContainer:    prodRepo.NewRepositoryContainer(pg),
	}
}
