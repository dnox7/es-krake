package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/pkg/utils"
)

const AttributeTableName = "attributes"

type AttributeRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) (entity.Attribute, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) ([]entity.Attribute, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Attribute, error)

	Update(ctx context.Context, attribute entity.Attribute, attributesToUpdate map[string]interface{}) (entity.Attribute, error)
}
