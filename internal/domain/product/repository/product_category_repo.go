package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type ProductCategoryRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.ProductCategory, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.ProductCategory, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributeValues []map[string]interface{}, batchSize int) error

	Update(ctx context.Context, prodCate entity.ProductCategory, attributesToUpdate map[string]interface{}) (entity.ProductCategory, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
