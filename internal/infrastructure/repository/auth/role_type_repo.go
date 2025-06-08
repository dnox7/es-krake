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

type roleTypeRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewRoleTypeRepository(pg *rdb.PostgreSQL) domainRepo.RoleTypeRepository {
	return &roleTypeRepo{
		logger: log.With("repository", "role_type_repo"),
		pg:     pg,
	}
}

// TakeByConditions implements repository.RoleTypeRepository.
func (r *roleTypeRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.RoleType, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return entity.RoleType{}, err
	}

	rt := entity.RoleType{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&rt).Error
	return rt, err
}

// FindByConditions implements repository.RoleTypeRepository.
func (r *roleTypeRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.RoleType, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	types := []entity.RoleType{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&types).Error
	return types, err
}
