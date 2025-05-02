package repository

import (
	"context"
	"fmt"
	"github.com/dpe27/es-krake/internal/domain/product/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"

	"gorm.io/gorm"
)

type attributeRepository struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewAttributeRepository(pg *rdb.PostgreSQL) domainRepo.AttributeRepository {
	return &attributeRepository{
		logger: log.With("repo", "attribute_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.AttributeRepository.
func (r *attributeRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Attribute, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
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
	spec specification.Base,
) ([]entity.Attribute, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
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
	err := utils.MapToStruct(attributes, &attributeEntity)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
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
		return entity.Attribute{}, fmt.Errorf(utils.ErrorGetTx)
	}

	attributeEntity := entity.Attribute{}
	err := utils.MapToStruct(attributes, &attributeEntity)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
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
	err := utils.MapToStruct(attributesToUpdate, &attribute)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
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
		return entity.Attribute{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &attribute)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Attribute{}, err
	}

	err = gormTx.Model(attribute).Updates(attributesToUpdate).Error
	return attribute, err
}
