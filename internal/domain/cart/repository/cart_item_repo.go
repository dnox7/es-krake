package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type CartItemRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.CartItem, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.CartItem, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.CartItem, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.CartItem, error)

	Update(ctx context.Context, cartItem entity.CartItem, attributesToUpdate map[string]interface{}) (entity.CartItem, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, cartItem entity.CartItem, attributesToUpdate map[string]interface{}) (entity.CartItem, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
