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

type funcCodeRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewFunctionCodeRepository(pg *rdb.PostgreSQL) domainRepo.FunctionCodeRepository {
	return &funcCodeRepo{
		logger: log.With("repository", "function_code_repo"),
		pg:     pg,
	}
}

// FindByConditions implements repository.FunctionCodeRepository.
func (f *funcCodeRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.FunctionCode, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		f.logger.Error(ctx, err.Error())
		return nil, err
	}

	funcCodes := []entity.FunctionCode{}
	err = f.pg.DB.
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&funcCodes).Error
	return funcCodes, err
}
