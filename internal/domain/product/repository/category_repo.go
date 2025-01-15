package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	"pech/es-krake/internal/domain/utils"
)

type CategoryRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.Category, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Category, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Category, error)

	FindByConditionsWithScope(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) ([]entity.Category, error)

	CreateWithTx(tx interface{}, attributes map[string]interface{}) (entity.Category, error)

	UpdateWithTx(tx interface{}, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)

	DeleteByConditionsWithTx(tx interface{}, conditions map[string]interface{}) error
}
