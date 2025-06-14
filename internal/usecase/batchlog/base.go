package usecase

import (
	"github.com/dpe27/es-krake/internal/domain/batchlog/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
)

type BatchLogUsecaseDeps struct {
	BatchLogRepo     repository.BatchLogRepository
	BatchLogTypeRepo repository.BatchLogTypeRepository
	Cache            redis.RedisRepository
}

type BatchLogUsecase struct {
	logger *log.Logger
	deps   *BatchLogUsecaseDeps
}

func NewBatchLogUsecase(deps *BatchLogUsecaseDeps) BatchLogUsecase {
	return BatchLogUsecase{
		logger: log.With("object", "batch_log_usecase"),
		deps:   deps,
	}
}
