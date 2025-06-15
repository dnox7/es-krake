package repository

import (
	"context"
	"errors"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type otpRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewOtpRepository(pg *rdb.PostgreSQL) domainRepo.OtpRepository {
	return &otpRepo{
		pg:     pg,
		logger: log.With("repository", "otp_repo"),
	}
}

// TakeByConditions implements repository.OtpRepository.
func (o *otpRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Otp, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		o.logger.Error(ctx, err.Error())
		return entity.Otp{}, err
	}

	otp := entity.Otp{}
	err = o.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&otp).Error
	return otp, err
}

// CreateWithTx implements repository.OtpRepository.
func (o *otpRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Otp, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Otp{}, errors.New(utils.ErrorGetTx)
	}

	otp := entity.Otp{}
	err := utils.MapToStruct(attributes, &otp)
	if err != nil {
		o.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Otp{}, err
	}

	err = gormTx.Create(&otp).Error
	return otp, err
}

// UpdateWithTx implements repository.OtpRepository.
func (o *otpRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	otp entity.Otp,
	attributesToUpdate map[string]interface{},
) (entity.Otp, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Otp{}, errors.New(utils.ErrorGetTx)
	}
	err := utils.MapToStruct(attributesToUpdate, &otp)
	if err != nil {
		o.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Otp{}, err
	}
	err = gormTx.Model(otp).Updates(attributesToUpdate).Error
	return otp, err
}

// DeleteByConditionsWithTx implements repository.OtpRepository.
func (o *otpRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return errors.New(utils.ErrorGetTx)
	}

	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		o.logger.Error(ctx, err.Error())
		return err
	}
	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.Otp{}).Error
}
