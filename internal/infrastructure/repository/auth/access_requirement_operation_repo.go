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

type accessReqOpRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewAccessRequirementOperationRepository(
	pg *rdb.PostgreSQL,
) domainRepo.AccessRequirementOperationRepository {
	return &accessReqOpRepo{
		logger: log.With("repository", "access_requirement_operation_repo"),
		pg:     pg,
	}
}

// PluckOperationIDByConditions implements repository.AccessRequirementOperationRepository.
func (a *accessReqOpRepo) PluckOperationIDByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]int, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return nil, err
	}

	var opIDs []int
	err = a.pg.GormDB().
		WithContext(ctx).
		Model(entity.AccessRequirementOperation{}).
		Scopes(scopes...).
		Where(conditions).
		Pluck("access_operation_id", &opIDs).
		Error
	return opIDs, err
}
