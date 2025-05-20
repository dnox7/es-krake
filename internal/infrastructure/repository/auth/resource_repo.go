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

type resourceRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewResourceRepository(pg *rdb.PostgreSQL) domainRepo.ResourceRepository {
	return &resourceRepo{
		logger: log.With("repo", "resource_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.ResourceRepository.
func (r *resourceRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Resource, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return entity.Resource{}, err
	}

	resoure := entity.Resource{}
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&resoure).Error
	return resoure, err
}

// FindByConditions implements repository.ResourceRepository.
func (r *resourceRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Resource, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	resoures := []entity.Resource{}
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&resoures).Error
	return resoures, err
}
