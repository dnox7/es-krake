package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/batchlog/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type BatchLogRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.BatchLog, error)
	CountByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (int64, error)
	Create(ctx context.Context, attributes map[string]interface{}) (entity.BatchLog, error)
	Update(cxt context.Context, batchLog entity.BatchLog, attributesToUpdate map[string]interface{}) (entity.BatchLog, error)
}
