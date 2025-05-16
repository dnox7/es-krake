package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const FunctionCodeTableName = "function_codes"

type FunctionCodeRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.FunctionCode, error)
}
