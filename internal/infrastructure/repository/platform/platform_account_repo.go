package repository

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type platformAccRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewPlatformAccountRepository(pg *rdb.PostgreSQL) domainRepo.PlatformAccountRepository {
	return &platformAccRepo{
		logger: log.With("repository", "platform_account_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.PlatformAccountRepository.
func (p *platformAccRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.PlatformAccount, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.PlatformAccount{}, err
	}

	var acc entity.PlatformAccount
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&acc).Error
	return acc, err
}

// FindByConditions implements repository.PlatformAccountRepository.
func (p *platformAccRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.PlatformAccount, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	accs := []entity.PlatformAccount{}
	err = p.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&accs).Error
	return accs, err
}

// Create implements repository.PlatformAccountRepository.
func (p *platformAccRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.PlatformAccount, error) {
	acc := entity.PlatformAccount{}
	err := utils.MapToStruct(attributes, &acc)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccount{}, err
	}

	err = p.pg.DB.WithContext(ctx).Create(&acc).Error
	return acc, err
}

// CreateWithTx implements repository.PlatformAccountRepository.
func (p *platformAccRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.PlatformAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.PlatformAccount{}, fmt.Errorf(utils.ErrorGetTx)
	}

	acc := entity.PlatformAccount{}
	err := utils.MapToStruct(attributes, &acc)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccount{}, err
	}

	err = gormTx.Create(&acc).Error
	return acc, err
}

// Update implements repository.PlatformAccountRepository.
func (p *platformAccRepo) Update(
	ctx context.Context,
	acc entity.PlatformAccount,
	attributesToUpdate map[string]interface{},
) (entity.PlatformAccount, error) {
	err := utils.MapToStruct(attributesToUpdate, &acc)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccount{}, err
	}

	err = p.pg.DB.
		WithContext(ctx).
		Model(acc).
		Updates(attributesToUpdate).Error
	return acc, err
}

// UpdateWithTx implements repository.PlatformAccountRepository.
func (p *platformAccRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	acc entity.PlatformAccount,
	attributesToUpdate map[string]interface{},
) (entity.PlatformAccount, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.PlatformAccount{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &acc)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.PlatformAccount{}, err
	}

	err = gormTx.Model(acc).Updates(attributesToUpdate).Error
	return acc, err
}

// DeleteByConditions implements repository.PlatformAccountRepository.
func (p *platformAccRepo) DeleteByConditions(
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
		p.logger.Error(ctx, err.Error())
		return err
	}
	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.PlatformAccount{}).Error
}
