package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/buyer/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type BuyerAccountRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.BuyerAccount, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.BuyerAccount, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.BuyerAccount, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.BuyerAccount, error)

	Update(ctx context.Context, buyerAccount entity.BuyerAccount, attributesToUpdate map[string]interface{}) (entity.BuyerAccount, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, buyerAccount entity.BuyerAccount, attributesToUpdate map[string]interface{}) (entity.BuyerAccount, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
