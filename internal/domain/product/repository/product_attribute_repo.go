package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const ProductAttributeValueTableName = "product_attribute_values"

type ProductAttributeValueRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductAttributeValue, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.ProductAttributeValue, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error)

	CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error)

	CreateBatch(ctx context.Context, attributeValues []map[string]interface{}, batchSize int) ([]entity.ProductAttributeValue, error)

	UpdateWithTx(ctx context.Context, attributeValue entity.ProductAttributeValue, attributesToUpdate map[string]interface{}) (entity.ProductAttributeValue, error)
}
