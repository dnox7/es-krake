package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

type platformAccountEnterpriseAccessRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewPlatformAccountEnterpriseAccessRepo(pg *rdb.PostgreSQL) domainRepo.PlatformAccountEnterpriseAccessRepository {
	return &platformAccountEnterpriseAccessRepo{
		pg:     pg,
		logger: log.With("repository", "platform_account_enterprise_access_repo"),
	}
}

// TakeByConditions implements repository.PlatformAccountEnterpriseAccessRepository.
func (p *platformAccountEnterpriseAccessRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.PlatformAccountEnterpriseAccess, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.PlatformAccountEnterpriseAccess{}, err
	}

	var pae entity.PlatformAccountEnterpriseAccess
	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&pae).Error
	return pae, err
}

// Create implements repository.PlatformAccountEnterpriseAccessRepository.
func (p *platformAccountEnterpriseAccessRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.PlatformAccountEnterpriseAccess, error) {
	pae := entity.PlatformAccountEnterpriseAccess{}
	err := utils.MapToStruct(attributes, &pae)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccountEnterpriseAccess{}, err
	}

	err = p.pg.GormDB().WithContext(ctx).Create(&pae).Error
	return pae, err
}

// Update implements repository.PlatformAccountEnterpriseAccessRepository.
func (p *platformAccountEnterpriseAccessRepo) Update(
	ctx context.Context,
	pae entity.PlatformAccountEnterpriseAccess,
	attributesToUpdate map[string]interface{},
) (entity.PlatformAccountEnterpriseAccess, error) {
	err := utils.MapToStruct(attributesToUpdate, &pae)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccountEnterpriseAccess{}, err
	}

	err = p.pg.GormDB().
		WithContext(ctx).
		Model(&pae).
		Updates(attributesToUpdate).
		Error
	return pae, err
}

// DeleteByConditions implements repository.PlatformAccountEnterpriseAccessRepository.
func (p *platformAccountEnterpriseAccessRepo) DeleteByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return err
	}

	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.PlatformAccountEnterpriseAccess{}).
		Error
	return err
}
