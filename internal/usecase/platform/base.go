package usecase

import (
	"github.com/dpe27/es-krake/internal/domain/platform/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
)

type PlatformUsecaseDeps struct {
	PlatformAccountRepo repository.PlatformAccountRepository
	DepartmentRepo      repository.DepartmentRepository

	Cache redis.RedisRepository
}

type PlatformUsecase struct {
	logger *log.Logger
	deps   *PlatformUsecaseDeps
}

func NewPlatformUsecase(deps *PlatformUsecaseDeps) PlatformUsecase {
	return PlatformUsecase{
		logger: log.With("object", "platform_usecase"),
		deps:   deps,
	}
}
