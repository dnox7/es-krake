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

type rolePermRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewRolePermissionRepository(pg *rdb.PostgreSQL) domainRepo.RolePermissionRepository {
	return &rolePermRepo{
		pg:     pg,
		logger: log.With("repository", "role_permission_repo"),
	}
}

// FindByConditions implements repository.RolePermissionRepository.
func (r *rolePermRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.RolePermission, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	rolePerms := []entity.RolePermission{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&rolePerms).Error
	return rolePerms, err
}

// CreateBatchWithTx implements repository.RolePermissionRepository.
func (r *rolePermRepo) CreateBatchWithTx(
	ctx context.Context,
	tx transaction.Base,
	attributes []map[string]interface{},
	batchSize int,
) error {
	gormTx, ok := tx.GetTx().(*gorm.DB)
	if !ok {
		return fmt.Errorf(utils.ErrorGetTx)
	}

	var (
		rp  entity.RolePermission
		err error
	)
	rolePerms := []entity.RolePermission{}
	for _, v := range attributes {
		err = utils.MapToStruct(v, &rp)
		if err != nil {
			r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
			return err
		}
		rolePerms = append(rolePerms, rp)
	}

	return gormTx.CreateInBatches(rolePerms, batchSize).Error
}

// DeleteByConditionsWithTx implements repository.RolePermissionRepository.
func (r *rolePermRepo) DeleteByConditionsWithTx(
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
		Delete(&entity.RolePermission{}).Error
}
