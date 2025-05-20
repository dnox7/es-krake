package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const ProductTableName = "products"

type ProductRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, scope specification.Base) (entity.Product, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, scope specification.Base) ([]entity.Product, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Product, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Product, error)

	Update(ctx context.Context, product entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, product entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error)
}
