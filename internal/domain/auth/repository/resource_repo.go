package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const ResourceTableName = "resoures"

type ResourceRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Resource, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Resource, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributes []map[string]interface{}, batchSize int) error
}
