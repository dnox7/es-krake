package repository

import (
	"context"
	"errors"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/cart/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type cartItemRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewCartItemRepository(pg *rdb.PostgreSQL) domainRepo.CartItemRepository {
	return &cartItemRepo{
		logger: log.With("repository", "cart_item_repo"),
		pg:     pg,
	}
}

// Create implements repository.CartItemRepository.
func (c *cartItemRepo) Create(ctx context.Context, attributes map[string]interface{}) (entity.CartItem, error) {
	cartItem := entity.CartItem{}
	err := utils.MapToStruct(attributes, &cartItem)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.CartItem{}, err
	}

	err = c.pg.GormDB().WithContext(ctx).Create(&cartItem).Error
	return cartItem, err
}

// CreateWithTx implements repository.CartItemRepository.
func (c *cartItemRepo) CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.CartItem, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.CartItem{}, errors.New(utils.ErrorGetTx)
	}

	cartItem := entity.CartItem{}
	err := utils.MapToStruct(attributes, &cartItem)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.CartItem{}, err
	}

	err = gormTx.Create(&cartItem).Error
	return cartItem, err
}

// DeleteByConditionsWithTx implements repository.CartItemRepository.
func (c *cartItemRepo) DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return errors.New(utils.ErrorGetTx)
	}

	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}

	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.CartItem{}).Error
}

// FindByConditions implements repository.CartItemRepository.
func (c *cartItemRepo) FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.CartItem, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return nil, err
	}

	cartItems := []entity.CartItem{}
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&cartItems).Error
	return cartItems, err
}

// TakeByConditions implements repository.CartItemRepository.
func (c *cartItemRepo) TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.CartItem, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return entity.CartItem{}, err
	}

	var cartItem entity.CartItem
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&cartItem).Error
	return cartItem, err
}

// Update implements repository.CartItemRepository.
func (c *cartItemRepo) Update(ctx context.Context, cartItem entity.CartItem, attributesToUpdate map[string]interface{}) (entity.CartItem, error) {
	err := utils.MapToStruct(attributesToUpdate, &cartItem)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.CartItem{}, err
	}

	err = c.pg.GormDB().WithContext(ctx).Model(&cartItem).Updates(attributesToUpdate).Error
	return cartItem, err
}

// UpdateWithTx implements repository.CartItemRepository.
func (c *cartItemRepo) UpdateWithTx(ctx context.Context, tx transaction.Base, cartItem entity.CartItem, attributesToUpdate map[string]interface{}) (entity.CartItem, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.CartItem{}, errors.New(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &cartItem)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.CartItem{}, err
	}

	err = gormTx.Model(&cartItem).Updates(attributesToUpdate).Error
	return cartItem, err
}
