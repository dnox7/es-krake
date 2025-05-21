package repository

import (
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	authRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/auth"
	enterpriseRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/enterprise"
	platformRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/platform"
	prodRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	AuthContainer       authRepo.RepositoryContainer
	EnterpriseContainer enterpriseRepo.RepositoryContainer
	PlatformContainer   platformRepo.RepositoryContainer
	ProductContainer    prodRepo.RepositoryContainer
}

func NewRepositoriesContainer(pg *rdb.PostgreSQL) *RepositoriesContainer {
	return &RepositoriesContainer{
		AuthContainer:       authRepo.NewRepositoryContainer(pg),
		EnterpriseContainer: enterpriseRepo.NewRepositoryContainer(pg),
		PlatformContainer:   platformRepo.NewRepositoryContainer(pg),
		ProductContainer:    prodRepo.NewRepositoryContainer(pg),
	}
}
