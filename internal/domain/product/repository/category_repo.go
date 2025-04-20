package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/domain/shared/transaction"
)

type CategoryRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) (entity.Category, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) ([]entity.Category, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Category, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Category, error)

	Update(ctx context.Context, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)
}
