package repository

import (
	"context"
	"fmt"

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

type roleRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewRoleRepository(pg *rdb.PostgreSQL) domainRepo.RoleRepository {
	return &roleRepo{
		pg:     pg,
		logger: log.With("repository", "role_repo"),
	}
}

// CheckExists implements repository.RoleRepository.
func (r *roleRepo) CheckExists(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (bool, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return false, err
	}

	var exists bool
	err = r.pg.GormDB().
		WithContext(ctx).
		Model(&entity.Role{}).
		Select("1").
		Where(conditions).
		Scopes(scopes...).
		Limit(1).
		Scan(&exists).Error
	return exists, err
}

// TakeByConditions implements repository.RoleRepository.
func (r *roleRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Role, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return entity.Role{}, err
	}

	role := entity.Role{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&role).Error
	return role, err
}

// FindByConditions implements repository.RoleRepository.
func (r *roleRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Role, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	roles := []entity.Role{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&roles).Error
	return roles, err
}

// CreateWithTx implements repository.RoleRepository.
func (r *roleRepo) CreateWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes map[string]interface{},
) (entity.Role, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Role{}, fmt.Errorf(utils.ErrorGetTx)
	}

	role := entity.Role{}
	err := utils.MapToStruct(attributes, &role)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Role{}, err
	}

	err = gormTx.Create(&role).Error
	return role, err
}

// UpdateWithTx implements repository.RoleRepository.
func (r *roleRepo) UpdateWithTx(
	ctx context.Context,
	tx transaction.Base,
	role entity.Role,
	attributesToUpdate map[string]interface{},
) (entity.Role, error) {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return entity.Role{}, fmt.Errorf(utils.ErrorGetTx)
	}

	err := utils.MapToStruct(attributesToUpdate, &role)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.Role{}, err
	}

	err = gormTx.Model(role).Updates(attributesToUpdate).Error
	return role, err
}

// DeleteByConditionsWithTx implements repository.RoleRepository.
func (r *roleRepo) DeleteByConditionsWithTx(
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
		r.logger.Error(ctx, err.Error())
		return err
	}
	return gormTx.
		Scopes(gormScopes...).
		Where(conditions).
		Delete(&entity.Role{}).Error
}
