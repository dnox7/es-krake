package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/seller/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type SellerAccountRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.SellerAccount, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.SellerAccount, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.SellerAccount, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.SellerAccount, error)

	Update(ctx context.Context, sellerAccount entity.SellerAccount, attributesToUpdate map[string]interface{}) (entity.SellerAccount, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, sellerAccount entity.SellerAccount, attributesToUpdate map[string]interface{}) (entity.SellerAccount, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
