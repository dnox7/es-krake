package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type PlatformAccountRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.PlatformAccount, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.PlatformAccount, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.PlatformAccount, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.PlatformAccount, error)

	Update(ctx context.Context, acc entity.PlatformAccount, attributesToUpdate map[string]interface{}) (entity.PlatformAccount, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, acc entity.PlatformAccount, attributesToUpdate map[string]interface{}) (entity.PlatformAccount, error)

	DeleteByConditions(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
