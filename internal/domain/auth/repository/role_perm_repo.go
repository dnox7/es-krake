package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type RolePermissionRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.RolePermission, error)

	CreateBatchWithTx(ctx context.Context, tx transaction.Base, attributes []map[string]interface{}, batchSize int) error

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
