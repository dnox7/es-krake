package repository

import (
	"context"
	"errors"

	"github.com/dpe27/es-krake/internal/domain/seller/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/seller/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type sellerAccountRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewSellerAccountRepository(pg *rdb.PostgreSQL) domainRepo.SellerAccountRepository {
	return &sellerAccountRepo{
		pg:     pg,
		logger: log.With("repository", "seller_account_repo"),
	}
}

// TakeByConditions implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.SellerAccount, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return entity.SellerAccount{}, err
	}

	var sellerAccount entity.SellerAccount
	err = s.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&sellerAccount).Error
	return sellerAccount, err
}

// FindByConditions implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.SellerAccount, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return []entity.SellerAccount{}, err
	}

	sellerAccounts := []entity.SellerAccount{}
	err = s.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&sellerAccounts).Error
	return sellerAccounts, err
}

// Create implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.SellerAccount, error) {
	sellerAccount := entity.SellerAccount{}
	err := utils.MapToStruct(attributes, &sellerAccount)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerAccount{}, err
	}

	err = s.pg.GormDB().WithContext(ctx).Create(&sellerAccount).Error
	return sellerAccount, err
}

// CreateWithTx implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.SellerAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.SellerAccount{}, errors.New(utils.ErrorGetTx)
	}

	sellerAccount := entity.SellerAccount{}
	err := utils.MapToStruct(attributes, &sellerAccount)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
	}

	err = gormTx.Create(&sellerAccount).Error
	return sellerAccount, err
}

// Update implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) Update(
	ctx context.Context,
	sellerAccount entity.SellerAccount,
	attributesToUpdate map[string]interface{},
) (entity.SellerAccount, error) {
	err := utils.MapToStruct(attributesToUpdate, &sellerAccount)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerAccount{}, err
	}

	err = s.pg.GormDB().
		WithContext(ctx).
		Model(&sellerAccount).
		Updates(attributesToUpdate).Error
	return sellerAccount, err
}

// UpdateWithTx implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	sellerAccount entity.SellerAccount,
	attributesToUpdate map[string]interface{},
) (entity.SellerAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.SellerAccount{}, errors.New(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &sellerAccount)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerAccount{}, err
	}

	err = gormTx.Model(&sellerAccount).Updates(attributesToUpdate).Error
	return sellerAccount, err
}

// DeleteByConditionsWithTx implements repository.SellerAccountRepository.
func (s *sellerAccountRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return err
	}

	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return errors.New(utils.ErrorGetTx)
	}

	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.SellerAccount{}).Error
}
