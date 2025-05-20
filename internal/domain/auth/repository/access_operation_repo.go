package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const AccessOperationsTableName = "access_operations"

type AccessOperationRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.AccessOperation, error)
}
