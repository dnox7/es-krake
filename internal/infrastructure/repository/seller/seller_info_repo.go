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

type sellerInfoRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewSellerInfoRepository(pg *rdb.PostgreSQL) domainRepo.SellerInfoRepository {
	return &sellerInfoRepo{
		pg:     pg,
		logger: log.With("repository", "seller_info_repo"),
	}
}

// TakeByConditions implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.SellerInfo, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return entity.SellerInfo{}, err
	}

	var sellerInfo entity.SellerInfo
	err = s.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&sellerInfo).Error
	return sellerInfo, err
}

// FindByConditions implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.SellerInfo, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return []entity.SellerInfo{}, err
	}

	sellerInfos := []entity.SellerInfo{}
	err = s.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&sellerInfos).Error
	return sellerInfos, err
}

// Create implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.SellerInfo, error) {
	sellerInfo := entity.SellerInfo{}
	err := utils.MapToStruct(attributes, &sellerInfo)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerInfo{}, err
	}

	err = s.pg.GormDB().WithContext(ctx).Create(&sellerInfo).Error
	return sellerInfo, err
}

// CreateWithTx implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.SellerInfo, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.SellerInfo{}, errors.New(utils.ErrorGetTx)
	}

	sellerInfo := entity.SellerInfo{}
	err := utils.MapToStruct(attributes, &sellerInfo)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerInfo{}, err
	}

	err = gormTx.Create(&sellerInfo).Error
	return sellerInfo, err
}

// Update implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) Update(
	ctx context.Context,
	sellerInfo entity.SellerInfo,
	attributesToUpdate map[string]interface{},
) (entity.SellerInfo, error) {
	err := utils.MapToStruct(attributesToUpdate, &sellerInfo)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerInfo{}, err
	}

	err = s.pg.GormDB().
		WithContext(ctx).
		Model(&sellerInfo).
		Updates(attributesToUpdate).Error
	return sellerInfo, err
}

// UpdateWithTx implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	sellerInfo entity.SellerInfo,
	attributesToUpdate map[string]interface{},
) (entity.SellerInfo, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.SellerInfo{}, errors.New(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &sellerInfo)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.SellerInfo{}, err
	}

	err = gormTx.Model(&sellerInfo).Updates(attributesToUpdate).Error
	return sellerInfo, err
}

// DeleteByConditionsWithTx implements repository.SellerInfoRepository.
func (s *sellerInfoRepo) DeleteByConditionsWithTx(
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
		Delete(&entity.SellerInfo{}).Error
}
