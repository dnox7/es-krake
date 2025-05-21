package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type departmentRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewDepartmentRepository(pg *rdb.PostgreSQL) domainRepo.DepartmentRepository {
	return &departmentRepo{
		logger: log.With("repository", "department_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.DepartmentRepository.
func (d *departmentRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Department, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		d.logger.Error(ctx, err.Error())
		return entity.Department{}, err
	}

	var dep entity.Department
	err = d.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&dep).Error
	return dep, err
}

// FindByConditions implements repository.DepartmentRepository.
func (d *departmentRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Department, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		d.logger.Error(ctx, err.Error())
		return nil, err
	}

	deps := []entity.Department{}
	err = d.pg.DB.
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&deps).Error
	return deps, err
}
