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

type enterpriseRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewEnterpriseRepository(
	pg *rdb.PostgreSQL,
) domainRepo.EnterpriseRepository {
	return &enterpriseRepo{
		logger: log.With("repository", "enterprise_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.EnterpriseRepository.
func (e *enterpriseRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Enterprise, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		e.logger.Error(ctx, err.Error())
		return entity.Enterprise{}, err
	}

	var ent entity.Enterprise
	err = e.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&ent).Error
	return ent, err
}

// FindByConditions implements repository.EnterpriseRepository.
func (e *enterpriseRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Enterprise, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		e.logger.Error(ctx, err.Error())
		return nil, err
	}

	ents := []entity.Enterprise{}
	err = e.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&ents).Error
	return ents, err
}

// PluckIDByConditions implements repository.EnterpriseRepository.
func (e *enterpriseRepo) PluckIDByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]int, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		e.logger.Error(ctx, err.Error())
		return nil, err
	}

	var IDs []int
	err = e.pg.DB.
		WithContext(ctx).
		Model(entity.Enterprise{}).
		Scopes(gormScopes...).
		Where(conditions).
		Pluck("id", &IDs).Error
	return IDs, err
}

// Create implements repository.EnterpriseRepository.
func (e *enterpriseRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.Enterprise, error) {
	ent := entity.Enterprise{}
	err := utils.MapToStruct(attributes, &ent)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Enterprise{}, err
	}

	err = e.pg.DB.WithContext(ctx).Create(&ent).Error
	return ent, err
}

// CreateWithTx implements repository.EnterpriseRepository.
func (e *enterpriseRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Enterprise, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Enterprise{}, fmt.Errorf(utils.ErrorGetTx)
	}

	ent := entity.Enterprise{}
	err := utils.MapToStruct(attributes, &ent)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Enterprise{}, err
	}

	err = gormTx.Create(&ent).Error
	return ent, err
}

// Update implements repository.EnterpriseRepository.
func (e *enterpriseRepo) Update(
	ctx context.Context,
	ent entity.Enterprise,
	attributesToUpdate map[string]interface{},
) (entity.Enterprise, error) {
	err := utils.MapToStruct(attributesToUpdate, &ent)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Enterprise{}, err
	}

	err = e.pg.DB.
		WithContext(ctx).
		Model(ent).
		Updates(attributesToUpdate).Error
	return ent, err
}

// UpdateWithTx implements repository.EnterpriseRepository.
func (e *enterpriseRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	ent entity.Enterprise,
	attributesToUpdate map[string]interface{},
) (entity.Enterprise, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Enterprise{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &ent)
	if err != nil {
		e.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Enterprise{}, err
	}

	err = gormTx.Model(ent).Updates(attributesToUpdate).Error
	return ent, err
}
