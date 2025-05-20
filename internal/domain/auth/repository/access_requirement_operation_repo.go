package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const AccessRequirementOperationTableName = "access_requirement_operations"

type AccessRequirementOperationRepository interface {
	PluckOperationIDByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]int, error)
}
