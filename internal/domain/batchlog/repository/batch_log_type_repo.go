package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/batchlog/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type BatchLogTypeRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.BatchLogType, error)
}
