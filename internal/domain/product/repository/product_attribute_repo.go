package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/domain/shared/transaction"
)

const ProductAttributeValueTableName = "product_attribute_values"

type ProductAttributeValueRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) (entity.ProductAttributeValue, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) ([]entity.ProductAttributeValue, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributeValues []map[string]interface{}, batchSize int) error

	Update(ctx context.Context, attributeValue entity.ProductAttributeValue, attributesToUpdate map[string]interface{}) (entity.ProductAttributeValue, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, scopes ...scope.Base) error
}
