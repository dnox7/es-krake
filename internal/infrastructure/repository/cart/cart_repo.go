package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/cart/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

type cartRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewCartRepository(pg *rdb.PostgreSQL) domainRepo.CartRepository {
	return &cartRepo{
		logger: log.With("repository", "cart_repo"),
		pg:     pg,
	}
}

// Create implements repository.CartRepository.
func (c *cartRepo) Create(ctx context.Context, attributes map[string]interface{}) (entity.Cart, error) {
	cart := entity.Cart{}
	err := utils.MapToStruct(attributes, &cart)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Cart{}, err
	}

	err = c.pg.GormDB().WithContext(ctx).Create(&cart).Error
	return cart, err
}

// FindByConditions implements repository.CartRepository.
func (c *cartRepo) FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Cart, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return nil, err
	}

	carts := []entity.Cart{}
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&carts).Error
	return carts, err
}

// TakeByConditions implements repository.CartRepository.
func (c *cartRepo) TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Cart, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return entity.Cart{}, err
	}

	var cart entity.Cart
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&cart).Error
	return cart, err
}

// Update implements repository.CartRepository.
func (c *cartRepo) Update(ctx context.Context, cart entity.Cart, attributesToUpdate map[string]interface{}) (entity.Cart, error) {
	err := utils.MapToStruct(attributesToUpdate, &cart)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Cart{}, err
	}

	err = c.pg.GormDB().
		WithContext(ctx).
		Model(&cart).
		Updates(attributesToUpdate).Error
	return cart, err
}
