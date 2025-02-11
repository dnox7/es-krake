package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const ProductOptionTableName = "product_options"

type IProductOptionRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductOption, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.ProductOption, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductOption, error)

	UpdateWithTx(ctx context.Context, option entity.ProductOption, attributesToUpdate map[string]interface{}) (entity.ProductOption, error)
}
