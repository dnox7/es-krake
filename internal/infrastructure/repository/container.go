package repository

import (
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	authRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/auth"
	batchLog "github.com/dpe27/es-krake/internal/infrastructure/repository/batchlog"
	enterpriseRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/enterprise"
	platformRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/platform"
	prodRepo "github.com/dpe27/es-krake/internal/infrastructure/repository/product"
)

type RepositoriesContainer struct {
	AuthContainer       authRepo.RepositoryContainer
	EnterpriseContainer enterpriseRepo.RepositoryContainer
	PlatformContainer   platformRepo.RepositoryContainer
	ProductContainer    prodRepo.RepositoryContainer
	BatchLogContainer   batchLog.RepositoryContainer
}

func NewRepositoriesContainer(pg *rdb.PostgreSQL, mongo *mdb.Mongo) *RepositoriesContainer {
	return &RepositoriesContainer{
		AuthContainer:       authRepo.NewRepositoryContainer(pg),
		EnterpriseContainer: enterpriseRepo.NewRepositoryContainer(pg),
		PlatformContainer:   platformRepo.NewRepositoryContainer(pg),
		ProductContainer:    prodRepo.NewRepositoryContainer(pg, mongo),
		BatchLogContainer:   batchLog.NewRepositoryContainer(pg),
	}
}
