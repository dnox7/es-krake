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

type resourceRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewResourceRepository(pg *rdb.PostgreSQL) domainRepo.ResourceRepository {
	return &resourceRepo{
		logger: log.With("repo", "resource_repo"),
		pg:     pg,
	}
}

// TakeByCondition implements repository.ResourceRepository.
func (r *resourceRepo) TakeByCondition(
	ctx context.Context,
	condition map[string]interface{},
	spec specification.Base,
) (entity.Resource, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return entity.Resource{}, err
	}

	resoure := entity.Resource{}
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(condition).
		Take(&resoure).Error
	return resoure, err
}

// FindByCondition implements repository.ResourceRepository.
func (r *resourceRepo) FindByCondition(
	ctx context.Context,
	condition map[string]interface{},
	spec specification.Base,
) ([]entity.Resource, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	resoures := []entity.Resource{}
	err = r.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(condition).
		Find(&resoures).Error
	return resoures, err
}

// CreateBatchWithTx implements repository.ResourceRepository.
func (r *resourceRepo) CreateBatchWithTx(
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
		rs  entity.Resource
		err error
	)
	resources := []entity.Resource{}
	for _, v := range attributes {
		err = utils.MapToStruct(v, &rs)
		if err != nil {
			r.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
			return err
		}
		resources = append(resources, rs)
	}

	return gormTx.CreateInBatches(resources, batchSize).Error
}
