package enterprise

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/enterprise/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type enterpriseAccRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewEnterpriseAccountRepository(
	pg *rdb.PostgreSQL,
) domainRepo.EnterpriseAccountRepository {
	return &enterpriseAccRepo{
		logger: log.With("repository", "enterprise_account_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.EnterpriseAccount, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		e.logger.Error(ctx, err.Error())
		return entity.EnterpriseAccount{}, err
	}

	var acc entity.EnterpriseAccount
	err = e.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&acc).Error
	return acc, err
}

// FindByConditions implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.EnterpriseAccount, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		e.logger.Error(ctx, err.Error())
		return nil, err
	}

	accs := []entity.EnterpriseAccount{}
	err = e.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&accs).Error
	return accs, err
}

// Create implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.EnterpriseAccount, error) {
	acc := entity.EnterpriseAccount{}
	err := utils.MapToStruct(attributes, &acc)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.EnterpriseAccount{}, err
	}

	err = e.pg.DB.WithContext(ctx).Create(&acc).Error
	return acc, err
}

// CreateWithTx implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.EnterpriseAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.EnterpriseAccount{}, fmt.Errorf(utils.ErrorGetTx)
	}

	acc := entity.EnterpriseAccount{}
	err := utils.MapToStruct(attributes, &acc)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.EnterpriseAccount{}, err
	}

	err = gormTx.Create(&acc).Error
	return acc, err
}

// Update implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) Update(
	ctx context.Context,
	acc entity.EnterpriseAccount,
	attributesToUpdate map[string]interface{},
) (entity.EnterpriseAccount, error) {
	err := utils.MapToStruct(attributesToUpdate, &acc)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.EnterpriseAccount{}, err
	}

	err = e.pg.DB.
		WithContext(ctx).
		Model(acc).
		Updates(attributesToUpdate).Error
	return acc, err
}

// UpdateWithTx implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	acc entity.EnterpriseAccount,
	attributesToUpdate map[string]interface{},
) (entity.EnterpriseAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.EnterpriseAccount{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &acc)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.EnterpriseAccount{}, err
	}

	err = gormTx.Model(acc).Updates(attributesToUpdate).Error
	return acc, err
}

// DeleteByConditionsWithTx implements repository.EnterpriseAccountRepository.
func (e *enterpriseAccRepo) DeleteByConditionsWithTx(
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
		e.logger.Error(ctx, err.Error())
		return err
	}
	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.EnterpriseAccount{}).Error
}
