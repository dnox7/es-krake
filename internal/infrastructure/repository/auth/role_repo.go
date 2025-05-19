package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/pkg/log"
)

type roleRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewRoleRepository(pg *rdb.PostgreSQL) domainRepo.RoleRepository {
	return &roleRepo{
		pg:     pg,
		logger: log.With("repository", "role_repo"),
	}
}

// TakeByConditions implements repository.RoleRepository.
func (r *roleRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Role, error) {
	panic("unimplemented")
}

// FindByConditions implements repository.RoleRepository.
func (r *roleRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Role, error) {
	panic("unimplemented")
}

// CreateWithTx implements repository.RoleRepository.
func (r *roleRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes []map[string]interface{},
) (entity.Role, error) {
	panic("unimplemented")
}

// UpdateWithTx implements repository.RoleRepository.
func (r *roleRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	role entity.Role,
	attributesToUpdate map[string]interface{},
) (entity.Role, error) {
	panic("unimplemented")
}

// DeleteByConditionsWithTx implements repository.RoleRepository.
func (r *roleRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	panic("unimplemented")
}
