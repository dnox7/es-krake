package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/specification"
	"pech/es-krake/internal/domain/shared/transaction"
)

const ProductOptionTableName = "product_options"

type ProductOptionRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.ProductOption, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.ProductOption, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributes []map[string]interface{}, batchSize int) error

	Update(ctx context.Context, option entity.ProductOption, attributesToUpdate map[string]interface{}) (entity.ProductOption, error)

	DeleteByConditionWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
