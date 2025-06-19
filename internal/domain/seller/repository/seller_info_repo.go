package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/seller/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type SellerInfoRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.SellerInfo, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.SellerInfo, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.SellerInfo, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.SellerInfo, error)

	Update(ctx context.Context, sellerInfo entity.SellerInfo, attributesToUpdate map[string]interface{}) (entity.SellerInfo, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, sellerInfo entity.SellerInfo, attributesToUpdate map[string]interface{}) (entity.SellerInfo, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
