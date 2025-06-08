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

type productRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewProductRepository(pg *rdb.PostgreSQL) domainRepo.ProductRepository {
	return &productRepo{
		logger: log.With("repository", "product_repo"),
		pg:     pg,
	}
}

// Create implements repository.ProductRepository.
func (p *productRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Product, error) {
	product := entity.Product{}
	err := utils.MapToStruct(attributes, &product)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Product{}, err
	}

	err = p.pg.GormDB().WithContext(ctx).Create(&product).Error
	return product, err
}

// CreateWithTx implements repository.ProductRepository.
func (p *productRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Product, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Product{}, fmt.Errorf(utils.ErrorGetTx)
	}

	product := entity.Product{}
	err := utils.MapToStruct(attributes, &product)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Product{}, err
	}

	err = gormTx.Create(&product).Error
	return product, err
}

// FindByConditions implements repository.ProductRepository.
func (p *productRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Product, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}
	products := []entity.Product{}
	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&products).Error
	return products, err
}

// TakeByConditions implements repository.ProductRepository.
func (p *productRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Product, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.Product{}, err
	}

	var product entity.Product
	err = p.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&product).Error
	return product, err
}

// Update implements repository.ProductRepository.
func (p *productRepo) Update(
	ctx context.Context,
	product entity.Product,
	attributesToUpdate map[string]interface{},
) (entity.Product, error) {
	err := utils.MapToStruct(attributesToUpdate, &product)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Product{}, err
	}

	err = p.pg.GormDB().
		WithContext(ctx).
		Model(product).
		Updates(attributesToUpdate).Error
	return product, err
}

// UpdateWithTx implements repository.ProductRepository.
func (p *productRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	product entity.Product,
	attributesToUpdate map[string]interface{},
) (entity.Product, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Product{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &product)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Product{}, err
	}

	err = gormTx.Model(product).Updates(attributesToUpdate).Error
	return product, err
}
