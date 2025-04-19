package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/shared/scope"
	"pech/es-krake/internal/domain/shared/transaction"
)

const AttributeTableName = "attributes"

type AttributeRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) (entity.Attribute, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...scope.Base) ([]entity.Attribute, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Attribute, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Attribute, error)

	Update(ctx context.Context, attribute entity.Attribute, attributesToUpdate map[string]interface{}) (entity.Attribute, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, attribute entity.Attribute, attributesToUpdate map[string]interface{}) (entity.Attribute, error)
}
