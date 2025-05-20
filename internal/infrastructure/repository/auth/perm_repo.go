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

type permRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewPermissionRepository(pg *rdb.PostgreSQL) domainRepo.PermissionRepository {
	return &permRepo{
		logger: log.With("repo", "permission_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.PermissionRepository.
func (p *permRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Permission, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.Permission{}, err
	}

	perm := entity.Permission{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&perm).Error
	return perm, err
}

// FindByConditions implements repository.PermissionRepository.
func (p *permRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Permission, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	perms := []entity.Permission{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&perms).Error
	return perms, err
}
