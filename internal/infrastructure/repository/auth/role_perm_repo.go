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

type rolePermRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewRolePermissionRepository(pg *rdb.PostgreSQL) domainRepo.RolePermissionRepository {
	return &rolePermRepo{
		pg:     pg,
		logger: log.With("repository", "role_permission_repo"),
	}
}

// FindByConditions implements repository.RolePermissionRepository.
func (r *rolePermRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.RolePermission, error) {
	panic("unimplemented")
}

// CreateBatchWithTx implements repository.RolePermissionRepository.
func (r *rolePermRepo) CreateBatchWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes []map[string]interface{},
	batchSize int,
) error {
	panic("unimplemented")
}

// CreateWithTx implements repository.RolePermissionRepository.
func (r *rolePermRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.RolePermission, error) {
	panic("unimplemented")
}

// DeleteByConditionsWithTx implements repository.RolePermissionRepository.
func (r *rolePermRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	panic("unimplemented")
}
