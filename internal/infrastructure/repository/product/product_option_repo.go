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

type productOptionRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewProductOptionRepository(pg *rdb.PostgreSQL) domainRepo.ProductOptionRepository {
	return &productOptionRepo{
		logger: log.With("repository", "product_option_repo"),
		pg:     pg,
	}
}

// CreateBatchWithTx implements repository.ProductOptionRepository.
func (p *productOptionRepo) CreateBatchWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes []map[string]interface{},
	batchSize int,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return fmt.Errorf(utils.ErrorGetTx)
	}

	var (
		opt entity.ProductOption
		err error
	)
	optSlice := []entity.ProductOption{}
	for _, v := range attributes {
		err = utils.MapToStruct(v, &opt)
		if err != nil {
			p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
			return err
		}
		optSlice = append(optSlice, opt)
	}

	return gormTx.CreateInBatches(optSlice, batchSize).Error
}

// DeleteByConditionWithTx implements repository.ProductOptionRepository.
func (p *productOptionRepo) DeleteByConditionWithTx(
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
		Delete(&entity.ProductOption{}).Error
}

// FindByConditions implements repository.ProductOptionRepository.
func (p *productOptionRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.ProductOption, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	optSlice := []entity.ProductOption{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&optSlice).Error
	return optSlice, err
}

// TakeByConditions implements repository.ProductOptionRepository.
func (p *productOptionRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.ProductOption, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.ProductOption{}, err
	}

	opt := entity.ProductOption{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&opt).Error
	return opt, err
}

// Update implements repository.ProductOptionRepository.
func (p *productOptionRepo) Update(
	ctx context.Context,
	option entity.ProductOption,
	attributesToUpdate map[string]interface{},
) (entity.ProductOption, error) {
	err := utils.MapToStruct(attributesToUpdate, &option)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.ProductOption{}, err
	}

	err = p.pg.DB.
		WithContext(ctx).
		Model(option).
		Updates(attributesToUpdate).Error
	return option, err
}
