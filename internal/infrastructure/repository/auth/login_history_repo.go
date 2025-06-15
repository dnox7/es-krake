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

type loginHistoryRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewLoginHistoryRepository(pg *rdb.PostgreSQL) domainRepo.LoginHistoryRepository {
	return &loginHistoryRepo{
		pg:     pg,
		logger: log.With("repository", "login_history_repo"),
	}
}

// TakeByConditions implements repository.LoginHistoryRepository.
func (l *loginHistoryRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.LoginHistory, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		l.logger.Error(ctx, err.Error())
		return entity.LoginHistory{}, err
	}

	loginHistory := entity.LoginHistory{}
	err = l.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&loginHistory).Error
	return loginHistory, err
}

// FindByConditions implements repository.LoginHistoryRepository.
func (l *loginHistoryRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.LoginHistory, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		l.logger.Error(ctx, err.Error())
		return nil, err
	}

	loginHistories := []entity.LoginHistory{}
	err = l.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&loginHistories).Error
	return loginHistories, err
}

// Create implements repository.LoginHistoryRepository.
func (l *loginHistoryRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.LoginHistory, error) {
	loginHistory := entity.LoginHistory{}
	err := utils.MapToStruct(attributes, &loginHistory)
	if err != nil {
		l.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.LoginHistory{}, err
	}

	err = l.pg.GormDB().WithContext(ctx).Create(&loginHistory).Error
	return loginHistory, err
}

// DeleteByConditionsWithTx implements repository.LoginHistoryRepository.
func (l *loginHistoryRepo) DeleteByConditionsWithTx(
	ctx context.Context,
	tx transaction.Base,
	conditions map[string]interface{},
	spec specification.Base,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		l.logger.Error(ctx, utils.ErrorGetTx)
		return errors.New(utils.ErrorGetTx)
	}

	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		l.logger.Error(ctx, err.Error())
		return err
	}

	err = gormTx.WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).Delete(&entity.LoginHistory{}).Error
	return err
}
