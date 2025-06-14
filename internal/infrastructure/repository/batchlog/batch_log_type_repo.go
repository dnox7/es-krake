package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/batchlog/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/batchlog/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type batchLogTypeRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewBatchLogTypeRepository(pg *rdb.PostgreSQL) domainRepo.BatchLogTypeRepository {
	return &batchLogTypeRepo{
		logger: log.With("repository", "batch_log_type_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.BatchLogTypeRepository.
func (b *batchLogTypeRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.BatchLogType, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		b.logger.Error(ctx, err.Error())
		return entity.BatchLogType{}, err
	}
	var batchLogType entity.BatchLogType
	err = b.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&batchLogType).Error
	return batchLogType, err
}
