package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type OptionAttributeValueRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.OptionAttributeValue, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.OptionAttributeValue, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributeValues []map[string]interface{}, batchSize int) error

	Update(ctx context.Context, attributeValue entity.OptionAttributeValue, attributesToUpdate map[string]interface{}) (entity.OptionAttributeValue, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
