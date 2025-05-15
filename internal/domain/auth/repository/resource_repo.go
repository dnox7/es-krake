package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const ResourceTableName = "resoures"

type ResourceRepository interface {
	TakeByCondition(ctx context.Context, condition map[string]interface{}, spec specification.Base) (entity.Resource, error)
	FindByCondition(ctx context.Context, condition map[string]interface{}, spec specification.Base) ([]entity.Resource, error)
	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributes []map[string]interface{}, batchSize int) error
}
