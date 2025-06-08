package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type permOperationRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewPermissionOperationRepository(
	pg *rdb.PostgreSQL,
) domainRepo.PermissionOperationRepository {
	return &permOperationRepo{
		logger: log.With("repository", "permission_operation_repo"),
		pg:     pg,
	}
}

// PluckOperationIDByConditions implements repository.PermissionOperationRepository.
func (p *permOperationRepo) PluckOperationIDByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]int, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	var opIDs []int
	err = p.pg.GormDB().
		WithContext(ctx).
		Model(entity.PermissionOperation{}).
		Scopes(scopes...).
		Where(conditions).
		Pluck("access_operation_id", &opIDs).
		Error
	return opIDs, err
}
