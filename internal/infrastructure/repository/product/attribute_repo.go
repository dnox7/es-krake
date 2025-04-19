package repository

import (
	"context"
	"fmt"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	domainScope "pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/domain/shared/transaction"
	"pech/es-krake/internal/infrastructure/db"
	gormScope "pech/es-krake/internal/infrastructure/db/gorm/scope"
	"pech/es-krake/pkg/log"
	baseUtils "pech/es-krake/pkg/utils"

	"gorm.io/gorm"
)

type attributeRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewAttributeRepository(pg *db.PostgreSQL) domainRepo.AttributeRepository {
	return &attributeRepository{
		logger: log.With("repo", "attribute_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.AttributeRepository.
func (r *attributeRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	scopes ...domainScope.Base,
) (entity.Attribute, error) {
	gormScopes, err := gormScope.ToGormScopes(scopes...)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return entity.Attribute{}, err
	}
	var attribute entity.Attribute
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&attribute).Error
	return attribute, err
}

// FindByConditions implements repository.AttributeRepository.
func (r *attributeRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	scopes ...domainScope.Base,
) ([]entity.Attribute, error) {
	gormScopes, err := gormScope.ToGormScopes(scopes...)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}
	attributes := []entity.Attribute{}
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&attributes).Error
	return attributes, err

}

// Create implements repository.AttributeRepository.
func (r *attributeRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Attribute, error) {
	attributeEntity := entity.Attribute{}
	err := baseUtils.MapToStruct(attributes, &attributeEntity)
	if err != nil {
		r.logger.Error(ctx, baseUtils.ErrorMapToStruct, "error", err.Error())
		return entity.Attribute{}, err
	}

	err = r.pg.DB.WithContext(ctx).Create(&attributeEntity).Error
	return attributeEntity, err
}

// CreateWithTx implements repository.AttributeRepository.
func (r *attributeRepository) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Attribute, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Attribute{}, fmt.Errorf(baseUtils.ErrorGetTx)
	}

	attributeEntity := entity.Attribute{}
	err := baseUtils.MapToStruct(attributes, &attributeEntity)
	if err != nil {
		r.logger.Error(ctx, baseUtils.ErrorMapToStruct, "error", err.Error())
		return entity.Attribute{}, err
	}

	err = gormTx.Create(&attributeEntity).Error
	return attributeEntity, err
}

// Update implements repository.AttributeRepository.
func (r *attributeRepository) Update(
	ctx context.Context,
	attribute entity.Attribute,
	attributesToUpdate map[string]interface{},
) (entity.Attribute, error) {
	err := baseUtils.MapToStruct(attributesToUpdate, &attribute)
	if err != nil {
		r.logger.Error(ctx, baseUtils.ErrorMapToStruct, "error", err.Error())
		return entity.Attribute{}, err
	}

	err = r.pg.DB.
		WithContext(ctx).
		Model(attribute).
		Updates(attributesToUpdate).Error
	return attribute, err
}

// UpdateWithTx implements repository.AttributeRepository.
func (r *attributeRepository) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attribute entity.Attribute,
	attributesToUpdate map[string]interface{},
) (entity.Attribute, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Attribute{}, fmt.Errorf(baseUtils.ErrorGetTx)
	}

	err := baseUtils.MapToStruct(attributesToUpdate, &attribute)
	if err != nil {
		r.logger.Error(ctx, baseUtils.ErrorMapToStruct, "error", err.Error())
		return entity.Attribute{}, err
	}

	err = gormTx.Model(attribute).Updates(attributesToUpdate).Error
	return attribute, err
}
