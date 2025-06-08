package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type PermissionOperationRepository interface {
	PluckOperationIDByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]int, error)
}
