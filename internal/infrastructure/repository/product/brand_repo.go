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

type brandRepository struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewBrandRepository(pg *rdb.PostgreSQL) domainRepo.BrandRepository {
	return &brandRepository{
		logger: log.With("repo", "brand_repo"),
		pg:     pg,
	}
}

// Create implements repository.BrandRepository.
func (b *brandRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Brand, error) {
	brand := entity.Brand{}
	err := utils.MapToStruct(attributes, &brand)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Brand{}, err
	}

	err = b.pg.DB.WithContext(ctx).Create(&brand).Error
	return brand, err
}

// FindByConditions implements repository.BrandRepository.
func (b *brandRepository) FindByConditions(
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
	err = b.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&brands).Error
	return brands, err
}

// TakeByConditions implements repository.BrandRepository.
func (b *brandRepository) TakeByConditions(
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
	err = b.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&brand).Error
	return brand, err
}

// Update implements repository.BrandRepository.
func (b *brandRepository) Update(
	ctx context.Context,
	brand entity.Brand,
	attributesToUpdate map[string]interface{},
) (entity.Brand, error) {
	err := utils.MapToStruct(attributesToUpdate, &brand)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Brand{}, err
	}

	err = b.pg.DB.
		WithContext(ctx).
		Model(brand).
		Updates(attributesToUpdate).Error
	return brand, err
}
