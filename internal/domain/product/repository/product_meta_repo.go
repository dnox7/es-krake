package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
)

const ProductMetaTableName = "productmetas"

type ProductMetaRespository interface {
	TakeByConditions(ctx context.Context, filter interface{}, scopes ...scope.Base) (entity.ProductMeta, error)

	FindByConditions(ctx context.Context, filter interface{}, scopes ...scope.Base) ([]entity.ProductMeta, error)
}
