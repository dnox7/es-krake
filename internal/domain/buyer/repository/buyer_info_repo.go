package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/buyer/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type BuyerInfoRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.BuyerInfo, error)

	TakeByCondition(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.BuyerInfo, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.BuyerInfo, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.BuyerInfo, error)

	Update(ctx context.Context, buyerInfo entity.BuyerInfo, attributesToUpdate map[string]interface{}) (entity.BuyerInfo, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, buyerInfo entity.BuyerInfo, attributesToUpdate map[string]interface{}) (entity.BuyerInfo, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
