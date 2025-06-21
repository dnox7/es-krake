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

type socialLoginProviderRepo struct {
	pg     *rdb.PostgreSQL
	logger *log.Logger
}

func NewSocialLoginProviderRepository(pg *rdb.PostgreSQL) domainRepo.SocialLoginProviderRepository {
	return &socialLoginProviderRepo{
		pg:     pg,
		logger: log.With("repository", "social_login_provider_repo"),
	}
}

func (r *socialLoginProviderRepo) FindByCondition(
	ctx context.Context,
	conditions map[string]interface{},
	spec specification.Base,
) ([]entity.SocialLoginProvider, error) {
	scopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	socialLoginProviders := []entity.SocialLoginProvider{}
	err = r.pg.GormDB().
		WithContext(ctx).
		Scopes(scopes...).
		Where(conditions).
		Find(&socialLoginProviders).Error
	return socialLoginProviders, err
}
