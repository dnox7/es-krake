package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
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

// TakeByCondition implements repository.PermissionRepository.
func (p *permRepo) TakeByCondition(
	ctx context.Context,
	condition map[string]interface{},
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
		Where(condition).
		Take(&perm).Error
	return perm, err
}

// FindByCondition implements repository.PermissionRepository.
func (p *permRepo) FindByCondition(
	ctx context.Context,
	condition map[string]interface{},
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
		Where(condition).
		Find(&perms).Error
	return perms, err
}

// Create implements repository.PermissionRepository.
func (p *permRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Permission, error) {
	perm := entity.Permission{}
	err := utils.MapToStruct(attributes, &perm)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Permission{}, err
	}

	err = p.pg.DB.WithContext(ctx).Create(&perm).Error
	return perm, err
}

// Update implements repository.PermissionRepository.
func (p *permRepo) Update(
	ctx context.Context,
	perm entity.Permission,
	attributesToUpdate map[string]interface{},
) (entity.Permission, error) {
	err := utils.MapToStruct(attributesToUpdate, &perm)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Permission{}, err
	}

	err = p.pg.DB.
		WithContext(ctx).
		Model(perm).
		Updates(attributesToUpdate).Error
	return perm, err
}
