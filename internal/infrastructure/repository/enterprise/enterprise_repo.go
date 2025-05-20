package enterprise

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/pkg/log"
)

type enterpriseRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewEnterpriseRepository(
	pg *rdb.PostgreSQL,
) domainRepo.EnterpriseRepository {
	return &enterpriseRepo{
		logger: log.With("repository", "enterprise_repo"),
		pg:     pg,
	}
}
