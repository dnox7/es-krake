package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const ProductAttributeValueTableName = "product_attribute_values"

type IProductAttributeValueRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductAttributeValue, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.ProductAttributeValue, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error)

	Update(ctx context.Context, attributeValue entity.ProductAttributeValue, attributesToUpdate map[string]interface{}) (entity.ProductAttributeValue, error)

	DeleteByConditions(ctx context.Context, conditions map[string]interface{}) error
}
