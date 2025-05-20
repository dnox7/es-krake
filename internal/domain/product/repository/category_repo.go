package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const CategoryTableName = "categories"

type CategoryRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Category, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Category, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Category, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Category, error)

	Update(ctx context.Context, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, category entity.Category, attributesToUpdate map[string]interface{}) (entity.Category, error)
}
