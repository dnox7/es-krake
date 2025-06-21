package repository

import (
	"context"
	"errors"

	"github.com/dpe27/es-krake/internal/domain/buyer/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/buyer/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type buyerAccountRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewBuyerAccountRepository(pg *rdb.PostgreSQL) domainRepo.BuyerAccountRepository {
	return &buyerAccountRepo{
		pg:     pg,
		logger: log.With("repository", "buyer_account_repo"),
	}
}

// TakeByConditions implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.BuyerAccount, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return entity.BuyerAccount{}, err
	}

	var buyerAccount entity.BuyerAccount
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&buyerAccount).Error
	return buyerAccount, err
}

// FindByConditions implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.BuyerAccount, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return nil, err
	}

	buyerAccounts := []entity.BuyerAccount{}
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&buyerAccounts).Error
	return buyerAccounts, err
}

// Create implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.BuyerAccount, error) {
	buyerAccount := entity.BuyerAccount{}
	err := utils.MapToStruct(attributes, &buyerAccount)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerAccount{}, err
	}

	err = b.pg.GormDB().WithContext(ctx).Create(&buyerAccount).Error
	return buyerAccount, err
}

// CreateWithTx implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.BuyerAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.BuyerAccount{}, errors.New(utils.ErrorGetTx)
	}

	buyerAccount := entity.BuyerAccount{}
	err := utils.MapToStruct(attributes, &buyerAccount)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerAccount{}, err
	}

	err = gormTx.Create(&buyerAccount).Error
	return buyerAccount, err
}

// Update implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) Update(
	ctx context.Context,
	buyerAccount entity.BuyerAccount,
	attributesToUpdate map[string]interface{},
) (entity.BuyerAccount, error) {
	err := utils.MapToStruct(attributesToUpdate, &buyerAccount)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerAccount{}, err
	}

	err = b.pg.GormDB().
		WithContext(ctx).
		Model(&buyerAccount).
		Updates(attributesToUpdate).Error
	return buyerAccount, err
}

// UpdateWithTx implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	buyerAccount entity.BuyerAccount,
	attributesToUpdate map[string]interface{},
) (entity.BuyerAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.BuyerAccount{}, errors.New(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &buyerAccount)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerAccount{}, err
	}

	err = gormTx.Model(&buyerAccount).Updates(attributesToUpdate).Error
	return buyerAccount, err
}

// DeleteByConditionsWithTx implements repository.BuyerAccountRepository.
func (b *buyerAccountRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return err
	}

	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return errors.New(utils.ErrorGetTx)
	}

	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.BuyerAccount{}).Error
}
