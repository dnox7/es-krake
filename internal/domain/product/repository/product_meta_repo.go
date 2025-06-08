package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type ProductMetaRespository interface {
	TakeByConditions(ctx context.Context, filter interface{}, spec specification.Base) (entity.ProductMeta, error)

	FindByConditions(ctx context.Context, filter interface{}, spec specification.Base) ([]entity.ProductMeta, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductMeta, error)

	UpdateByID(ctx context.Context, ID interface{}, operation interface{}) (entity.ProductMeta, error)
}
