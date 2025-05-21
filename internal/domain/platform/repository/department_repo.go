package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type DepartmentRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Department, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Department, error)
}
