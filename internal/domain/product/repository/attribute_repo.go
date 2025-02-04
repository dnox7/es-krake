package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const AttributeTableName = "attributes"

type IAttributeRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.Attribute, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Attribute, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Attribute, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Attribute, error)

	Update(ctx context.Context, attribute entity.Attribute, attributesToUpdate map[string]interface{}) (entity.Attribute, error)
}
