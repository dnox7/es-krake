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

type productCategoryRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewProductCategoryRepository(pg *rdb.PostgreSQL) domainRepo.ProductCategoryRepository {
	return &productCategoryRepo{
		logger: log.With("repository", "product_category_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.ProductCategoryRepository.
func (p *productCategoryRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.ProductCategory, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.ProductCategory{}, err
	}
	pc := entity.ProductCategory{}
	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&pc).Error
	return pc, err
}

// FindByConditions implements repository.ProductCategoryRepository.
func (p *productCategoryRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.ProductCategory, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}
	pcs := []entity.ProductCategory{}
	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&pcs).Error
	return pcs, err
}

// CreateBatchWithTx implements repository.ProductCategoryRepository.
func (p *productCategoryRepo) CreateBatchWithTx(
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
		pc  entity.ProductCategory
		err error
	)
	pcs := []entity.ProductCategory{}
	for _, v := range attributeValues {
		err = utils.MapToStruct(v, &pc)
		if err != nil {
			p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
			return err
		}
		pcs = append(pcs, pc)
	}

	return gormTx.CreateInBatches(pcs, batchSize).Error
}

// Update implements repository.ProductCategoryRepository.
func (p *productCategoryRepo) Update(
	ctx context.Context,
	prodCate entity.ProductCategory,
	attributesToUpdate map[string]interface{},
) (entity.ProductCategory, error) {
	err := utils.MapToStruct(attributesToUpdate, &prodCate)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.ProductCategory{}, err
	}

	err = p.pg.GormDB().
		WithContext(ctx).
		Model(prodCate).
		Updates(attributesToUpdate).Error
	return prodCate, err
}

// DeleteByConditionsWithTx implements repository.ProductCategoryRepository.
func (p *productCategoryRepo) DeleteByConditionsWithTx(
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
		Delete(&entity.ProductCategory{}).Error
}
