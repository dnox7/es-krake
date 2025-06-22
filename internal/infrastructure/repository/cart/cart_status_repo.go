package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/cart/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	gormScope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type cartStatusRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewCartStatusRepository(pg *rdb.PostgreSQL) domainRepo.CartStatusRepository {
	return &cartStatusRepo{
		logger: log.With("repository", "cart_status_repo"),
		pg:     pg,
	}
}

// FindByConditions implements repository.CartStatusRepository.
func (c *cartStatusRepo) FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.CartStatus, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return nil, err
	}

	cartStatuses := []entity.CartStatus{}
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Find(&cartStatuses).Error
	return cartStatuses, err
}

// TakeByConditions implements repository.CartStatusRepository.
func (c *cartStatusRepo) TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.CartStatus, error) {
	gormScopes, err := gormScope.ToGormScopes(spec)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return entity.CartStatus{}, err
	}

	var cartStatus entity.CartStatus
	err = c.pg.GormDB().
		WithContext(ctx).
		Scopes(gormScopes...).
		Where(conditions).
		Take(&cartStatus).Error
	return cartStatus, err
}
