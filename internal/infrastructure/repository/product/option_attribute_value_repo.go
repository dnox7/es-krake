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

type optionAttributeValueRepository struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewOptionAttributeValueRepository(pg *rdb.PostgreSQL) domainRepo.OptionAttributeValueRepository {
	return &optionAttributeValueRepository{
		logger: log.With("repo", "product_attribute_value_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.OptionAttributeValue, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.OptionAttributeValue{}, err
	}
	pav := entity.OptionAttributeValue{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&pav).Error
	return pav, err
}

// FindByConditions implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.OptionAttributeValue, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}
	pavSlice := []entity.OptionAttributeValue{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&pavSlice).Error
	return pavSlice, err
}

// CreateBatchWithTx implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepository) CreateBatchWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributeValues []map[string]interface{},
	batchSize int,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return fmt.Errorf(utils.ErrorGetTx)
	}

	var (
		pav entity.OptionAttributeValue
		err error
	)
	pavSlice := []entity.OptionAttributeValue{}
	for _, v := range attributeValues {
		err = utils.MapToStruct(v, &pav)
		if err != nil {
			p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
			return err
		}
		pavSlice = append(pavSlice, pav)
	}

	return gormTx.CreateInBatches(pavSlice, batchSize).Error
}

// Update implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepository) Update(
	ctx context.Context,
	attributeValue entity.OptionAttributeValue,
	attributesToUpdate map[string]interface{},
) (entity.OptionAttributeValue, error) {
	err := utils.MapToStruct(attributesToUpdate, &attributeValue)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.OptionAttributeValue{}, err
	}

	err = p.pg.DB.
		WithContext(ctx).
		Model(attributeValue).
		Updates(attributesToUpdate).Error
	return attributeValue, err
}

// DeleteByConditions implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepository) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return fmt.Errorf(utils.ErrorGetTx)
	}

	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return err
	}
	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.OptionAttributeValue{}).Error
}
