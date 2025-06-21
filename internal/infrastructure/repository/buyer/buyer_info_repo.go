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

type buyerInfoRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

// DeleteByConditionsWithTx implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error {
	panic("unimplemented")
}

func NewBuyerInfoRepository(pg *rdb.PostgreSQL) domainRepo.BuyerInfoRepository {
	return &buyerInfoRepo{
		pg:     pg,
		logger: log.With("repository", "buyer_info_repo"),
	}
}

// TakeByCondition implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) TakeByCondition(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.BuyerInfo, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return entity.BuyerInfo{}, err
	}

	var buyerInfo entity.BuyerInfo
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&buyerInfo).Error
	return buyerInfo, err
}

// FindByConditions implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.BuyerInfo, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return nil, err
	}

	buyerInfos := []entity.BuyerInfo{}
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&buyerInfos).Error
	return buyerInfos, err
}

// Create implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.BuyerInfo, error) {
	buyerInfo := entity.BuyerInfo{}
	err := utils.MapToStruct(attributes, &buyerInfo)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerInfo{}, err
	}

	err = b.pg.GormDB().WithContext(ctx).Create(&buyerInfo).Error
	return buyerInfo, err
}

// CreateWithTx implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.BuyerInfo, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.BuyerInfo{}, errors.New(utils.ErrorGetTx)
	}

	buyerInfo := entity.BuyerInfo{}
	err := utils.MapToStruct(attributes, &buyerInfo)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerInfo{}, err
	}

	err = gormTx.Create(&buyerInfo).Error
	return buyerInfo, err
}

// Update implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) Update(
	ctx context.Context,
	buyerInfo entity.BuyerInfo,
	attributesToUpdate map[string]interface{},
) (entity.BuyerInfo, error) {
	err := utils.MapToStruct(attributesToUpdate, &buyerInfo)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerInfo{}, err
	}

	err = b.pg.GormDB().
		WithContext(ctx).
		Model(&buyerInfo).
		Updates(attributesToUpdate).Error
	return buyerInfo, err
}

// UpdateWithTx implements repository.BuyerInfoRepository.
func (b *buyerInfoRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	buyerInfo entity.BuyerInfo,
	attributesToUpdate map[string]interface{},
) (entity.BuyerInfo, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.BuyerInfo{}, errors.New(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &buyerInfo)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BuyerInfo{}, err
	}

	err = gormTx.Model(&buyerInfo).Updates(attributesToUpdate).Error
	return buyerInfo, err
}
