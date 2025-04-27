package repository

import (
	"context"
	"fmt"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/domain/shared/specification"
	"pech/es-krake/internal/domain/shared/transaction"
	"pech/es-krake/internal/infrastructure/rdb"
	gormScope "pech/es-krake/internal/infrastructure/rdb/gorm/scope"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	"gorm.io/gorm"
)

type categoryRepository struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewCategoryRepository(pg *rdb.PostgreSQL) domainRepo.CategoryRepository {
	return &categoryRepository{
		logger: log.With("repo", "category_repo"),
		pg:     pg,
	}
}

// Create implements repository.CategoryRepository.
func (c *categoryRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Category, error) {
	category := entity.Category{}
	err := utils.MapToStruct(attributes, &category)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Category{}, err
	}

	err = c.pg.DB.WithContext(ctx).Create(&category).Error
	return category, err
}

// CreateWithTx implements repository.CategoryRepository.
func (c *categoryRepository) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Category, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Category{}, fmt.Errorf(utils.ErrorGetTx)
	}

	category := entity.Category{}
	err := utils.MapToStruct(attributes, &category)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Category{}, err
	}

	err = gormTx.Create(&category).Error
	return category, err
}

// FindByConditions implements repository.CategoryRepository.
func (c *categoryRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Category, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return nil, err
	}

	categories := []entity.Category{}
	err = c.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&categories).Error
	return categories, err

}

// TakeByConditions implements repository.CategoryRepository.
func (c *categoryRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Category, error) {

	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return entity.Category{}, err
	}

	category := entity.Category{}
	err = c.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&category).Error
	return category, err
}

// Update implements repository.CategoryRepository.
func (c *categoryRepository) Update(
	ctx context.Context,
	category entity.Category,
	attributesToUpdate map[string]interface{},
) (entity.Category, error) {
	err := utils.MapToStruct(attributesToUpdate, &category)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Category{}, err
	}

	err = c.pg.DB.
		WithContext(ctx).
		Model(category).
		Updates(attributesToUpdate).Error
	return category, err
}

// UpdateWithTx implements repository.CategoryRepository.
func (c *categoryRepository) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	category entity.Category,
	attributesToUpdate map[string]interface{},
) (entity.Category, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Category{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &category)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Category{}, err
	}

	err = gormTx.Model(category).Updates(attributesToUpdate).Error
	return category, err
}
