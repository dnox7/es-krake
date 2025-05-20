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

type optionAttributeValueRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewOptionAttributeValueRepository(pg *rdb.PostgreSQL) domainRepo.OptionAttributeValueRepository {
	return &optionAttributeValueRepo{
		logger: log.With("repository", "product_attribute_value_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.OptionAttributeValueRepository.
func (p *optionAttributeValueRepo) TakeByConditions(
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
func (p *optionAttributeValueRepo) FindByConditions(
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
func (p *optionAttributeValueRepo) CreateBatchWithTx(
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
func (p *optionAttributeValueRepo) Update(
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
func (p *optionAttributeValueRepo) DeleteByConditionsWithTx(
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
