package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/cart/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type CartStatusType int

const (
	CartStatusTypeActive CartStatusType = iota + 1
	CartStatusTypeInactive
)

type CartStatusRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.CartStatus, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.CartStatus, error)
}
