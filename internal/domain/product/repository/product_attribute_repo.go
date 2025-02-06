package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const ProductAttributeValueTableName = "product_attribute_values"

type IProductAttributeValueRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.ProductAttributeValue, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductAttributeValue, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.ProductAttributeValue, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error)

	Update(ctx context.Context, attributeValue entity.ProductAttributeValue, attributesToUpdate map[string]interface{}) (entity.ProductAttributeValue, error)

	DeleteByID(ctx context.Context, ID int) error

	DeleteByConditions(ctx context.Context, conditions map[string]interface{}) error
}
