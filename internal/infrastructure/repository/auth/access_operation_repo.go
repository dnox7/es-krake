package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type accessOpetationRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewAccessOperationRepository(pg *rdb.PostgreSQL) domainRepo.AccessOperationRepository {
	return &accessOpetationRepo{
		logger: log.With("repository", "access_operation_repo"),
		pg:     pg,
	}
}

// FindByConditions implements repository.AccessOperationRepository.
func (f *accessOpetationRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.AccessOperation, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		f.logger.Error(ctx, err.Error())
		return nil, err
	}

	funcCodes := []entity.AccessOperation{}
	err = f.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&funcCodes).Error
	return funcCodes, err
}
