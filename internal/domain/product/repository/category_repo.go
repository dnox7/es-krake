package repository

import (
	"context"
	"pech/es-krake/internal/domain"
	"pech/es-krake/internal/domain/product/entity"
)

type ICategoryRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.Category, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Category, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Category, error)

	FindByConditionsWithScope(ctx context.Context, conditions map[string]interface{}, scopes ...domain.Scope) ([]entity.Category, error)

	CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Category, error)

	UpdateWithTx(ctx context.Context, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)
}
