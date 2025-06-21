package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type LoginHistoryRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.LoginHistory, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.LoginHistory, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.LoginHistory, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
