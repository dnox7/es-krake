package enterprise

import (
	domainRepo "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/pkg/log"
)

type enterpriseAccRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewEnterpriseAccountRepository(
	pg *rdb.PostgreSQL,
) domainRepo.EnterpriseAccountRepository {
	return &enterpriseAccRepo{
		logger: log.With("repository", "enterprise_account_repo"),
		pg:     pg,
	}
}
