package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

type brandRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewBrandRepository(pg *rdb.PostgreSQL) domainRepo.BrandRepository {
	return &brandRepo{
		logger: log.With("repository", "brand_repo"),
		pg:     pg,
	}
}

// Create implements repository.BrandRepository.
func (b *brandRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Brand, error) {
	brand := entity.Brand{}
	err := utils.MapToStruct(attributes, &brand)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Brand{}, err
	}

	err = b.pg.GormDB().WithContext(ctx).Create(&brand).Error
	return brand, err
}

// FindByConditions implements repository.BrandRepository.
func (b *brandRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Brand, error) {
	gormScopes, err := scope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return nil, err
	}

	brands := []entity.Brand{}
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&brands).Error
	return brands, err
}

// TakeByConditions implements repository.BrandRepository.
func (b *brandRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Brand, error) {
	gormScopes, err := scope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return entity.Brand{}, err
	}

	brand := entity.Brand{}
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&brand).Error
	return brand, err
}

// Update implements repository.BrandRepository.
func (b *brandRepo) Update(
	ctx context.Context,
	brand entity.Brand,
	attributesToUpdate map[string]interface{},
) (entity.Brand, error) {
	err := utils.MapToStruct(attributesToUpdate, &brand)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Brand{}, err
	}

	err = b.pg.GormDB().
		WithContext(ctx).
		Model(brand).
		Updates(attributesToUpdate).Error
	return brand, err
}
