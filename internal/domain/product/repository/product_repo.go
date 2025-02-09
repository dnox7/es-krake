package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const ProductTableName = "products"

type IProductRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.Product, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Product, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Product, error)

	Create(ctx context.Context, product map[string]interface{}) (entity.Product, error)

	Update(ctx context.Context, product entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error)
}
