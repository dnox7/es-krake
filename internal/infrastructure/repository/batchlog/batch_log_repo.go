package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/batchlog/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/batchlog/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

type batchLogRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewBatchLogRepository(pg *rdb.PostgreSQL) domainRepo.BatchLogRepository {
	return &batchLogRepo{
		logger: log.With("repository", "batch_log_repo"),
		pg:     pg,
	}
}

// CountByConditions implements repository.BatchLogRepository.
func (b *batchLogRepo) CountByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (int64, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return 0, err
	}

	var count int64
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Model(&entity.BatchLog{}).
		Where(conditions).
		Count(&count).Error
	return count, err
}

// Create implements repository.BatchLogRepository.
func (b *batchLogRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.BatchLog, error) {
	batchLog := entity.BatchLog{}
	err := utils.MapToStruct(attributes, &batchLog)
	if err != nil {
		b.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BatchLog{}, err
	}

	err = b.pg.GormDB().WithContext(ctx).Create(&batchLog).Error
	return batchLog, err
}

// TakeByConditions implements repository.BatchLogRepository.
func (b *batchLogRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.BatchLog, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return entity.BatchLog{}, err
	}
	var batchLog entity.BatchLog
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&batchLog).Error
	return batchLog, err
}

// Update implements repository.BatchLogRepository.
func (b *batchLogRepo) Update(
	cxt context.Context,
	batchLog entity.BatchLog,
	attributesToUpdate map[string]interface{},
) (entity.BatchLog, error) {
	err := utils.MapToStruct(attributesToUpdate, &batchLog)
	if err != nil {
		b.logger.Error(cxt, utils.ErrorMapToStruct, "error", err.Error())
		return entity.BatchLog{}, err
	}

	err = b.pg.GormDB().
		WithContext(cxt).
		Model(&batchLog).
		Updates(attributesToUpdate).Error
	return batchLog, err
}
