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

type actionRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewActionRepository(pg *rdb.PostgreSQL) domainRepo.ActionRepository {
	return &actionRepo{
		logger: log.With("repo", "action_repo"),
		pg:     pg,
	}
}

// TakeByCondition implements repository.ActionRepository.
func (a *actionRepo) TakeByCondition(
	ctx context.Context,
	condition map[string]interface{},
	spec specification.Base,
) (entity.Action, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return entity.Action{}, err
	}

	action := entity.Action{}
	err = a.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(condition).
		Take(&action).Error
	return action, err
}

// FindByCondition implements repository.ActionRepository.
func (a *actionRepo) FindByCondition(
	ctx context.Context,
	condition map[string]interface{},
	spec specification.Base,
) ([]entity.Action, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return nil, err
	}

	actions := []entity.Action{}
	err = a.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(condition).
		Find(&actions).Error
	return actions, err
}
