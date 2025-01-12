package repository

import (
	"context"
	"pech/es-krake/internal/domain/category/entity"
)

type CategoryAttributeRepository interface {
	FindByCategoryID(ctx context.Context, categoryID int) ([]entity.CategoryAttribute, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.CategoryAttribute, error)
}
