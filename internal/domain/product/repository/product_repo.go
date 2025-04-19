package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
)

const ProductTableName = "products"

type ProductRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) (entity.Product, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) ([]entity.Product, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Product, error)

	CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Product, error)

	UpdateWithTx(ctx context.Context, product entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error)
}
