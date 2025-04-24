package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/specification"
)

const ProductMetaTableName = "productmetas"

type ProductMetaRespository interface {
	TakeByConditions(ctx context.Context, filter interface{}, spec specification.Base) (entity.ProductMeta, error)

	FindByConditions(ctx context.Context, filter interface{}, spec specification.Base) ([]entity.ProductMeta, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductMeta, error)
}
