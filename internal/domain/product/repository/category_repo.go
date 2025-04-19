package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
)

type CategoryRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Category, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Category, error)

	FindByConditionsWithScope(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) ([]entity.Category, error)

	CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Category, error)

	UpdateWithTx(ctx context.Context, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)
}
