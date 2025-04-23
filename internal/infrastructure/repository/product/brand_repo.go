package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	domainScope "pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/infrastructure/db"
	gormScope "pech/es-krake/internal/infrastructure/db/gorm/scope"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"
)

type brandRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewBrandRepository(pg *db.PostgreSQL) domainRepo.BrandRepository {
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
	scopes ...domainScope.Base,
) ([]entity.Brand, error) {
	gormScopes, err := gormScope.ToGormScopes(scopes...)
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
	scopes ...domainScope.Base,
) (entity.Brand, error) {
	gormScopes, err := gormScope.ToGormScopes(scopes...)
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
