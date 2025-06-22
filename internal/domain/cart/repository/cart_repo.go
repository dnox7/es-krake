package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type CartRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Cart, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Cart, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Cart, error)

	Update(ctx context.Context, cart entity.Cart, attributesToUpdate map[string]interface{}) (entity.Cart, error)
}
