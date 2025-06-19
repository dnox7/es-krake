package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/buyer/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/buyer/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type genderRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewGenderRepository(pg *rdb.PostgreSQL) domainRepo.GenderRepository {
	return &genderRepo{
		pg:     pg,
		logger: log.With("repository", "gender_repo"),
	}
}

// TakeByCondition implements repository.GenderRepository.
func (g *genderRepo) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) (entity.Gender, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		g.logger.Error(ctx, err.Error())
		return entity.Gender{}, err
	}

	var gender entity.Gender
	err = g.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Take(&gender).Error
	return gender, err
}

// FindByConditions implements repository.GenderRepository.
func (g *genderRepo) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.Gender, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		g.logger.Error(ctx, err.Error())
		return nil, err
	}

	genders := []entity.Gender{}
	err = g.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&genders).Error
	return genders, err
}
