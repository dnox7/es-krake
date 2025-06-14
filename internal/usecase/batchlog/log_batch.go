package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dpe27/es-krake/internal/domain/batchlog/entity"
	"github.com/dpe27/es-krake/internal/interfaces/batch/dto"
	"gorm.io/gorm"
)

func (blu *BatchLogUsecase) LogBatchStarting(
	ctx context.Context,
	event dto.BatchEvent,
) (entity.BatchLog, error) {
	logType, err := blu.deps.BatchLogTypeRepo.TakeByConditions(ctx, map[string]interface{}{
		"type": event.Type,
	}, nil)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.BatchLog{}, fmt.Errorf(
				"the event type '%v' does not exists in the 'batch_log_types' table",
				event.Type,
			)
		}
		return entity.BatchLog{}, err
	}

	dataJson, err := json.Marshal(event.Data)
	if err != nil {
		return entity.BatchLog{}, err
	}

	return blu.deps.BatchLogRepo.Create(ctx, map[string]interface{}{
		"batch_log_type_id": logType.ID,
		"event_id":          event.ID,
		"arguments":         string(dataJson),
		"started_at":        time.Now(),
		"ended_at":          nil,
		"success":           nil,
		"error_message":     nil,
	})
}

func (blu *BatchLogUsecase) LogBatchEnded(
	ctx context.Context,
	batchLog entity.BatchLog,
	returnedErr error,
) (entity.BatchLog, error) {
	if returnedErr == nil {
		return blu.deps.BatchLogRepo.Update(ctx, batchLog, map[string]interface{}{
			"ended_at": time.Now(),
			"success":  false,
		})
	}

	return blu.deps.BatchLogRepo.Update(ctx, batchLog, map[string]interface{}{
		"ended_at":      time.Now(),
		"success":       true,
		"error_message": returnedErr.Error(),
	})
}
