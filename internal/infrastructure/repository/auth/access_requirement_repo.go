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

type accessRequirementRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewAccessRequirementRepository(pg *rdb.PostgreSQL) domainRepo.AccessRequirementRepository {
	return &accessRequirementRepo{
		logger: log.With("repository", "access_requirement_repo"),
	}
}

// TakeByConditions implements repository.AccessRequirementRepository.
func (a *accessRequirementRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.AccessRequirement, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return entity.AccessRequirement{}, err
	}

	accessReq := entity.AccessRequirement{}
	err = a.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&accessReq).Error
	return accessReq, err
}
