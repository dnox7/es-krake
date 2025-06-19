package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/buyer/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type GenderRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Gender, error)

	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Gender, error)
}
