package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/pkg/utils"
)

const ProductOptionTableName = "product_options"

type ProductOptionRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) (entity.ProductOption, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) ([]entity.ProductOption, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductOption, error)

	UpdateWithTx(ctx context.Context, option entity.ProductOption, attributesToUpdate map[string]interface{}) (entity.ProductOption, error)
}
