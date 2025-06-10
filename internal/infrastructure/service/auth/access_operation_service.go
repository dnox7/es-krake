package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	goredis "github.com/redis/go-redis/v9"

	// "github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
)

type accessOperationService struct {
	accessOpRepo repository.AccessOperationRepository
	cache        redis.RedisRepository
	logger       *log.Logger
}

func NewAccessOperationService(
	accessOpRepo repository.AccessOperationRepository,
	cache redis.RedisRepository,
) domainService.AccessOperationService {
	return &accessOperationService{
		accessOpRepo,
		cache,
		log.With("service", "access_operation_service"),
	}
}

// GetOperationsWithAccessReqCode implements service.AccessOperationService.
func (a *accessOperationService) GetOperationsWithAccessReqCode(
	ctx context.Context,
	code string,
) ([]entity.AccessOperation, error) {
	cacheKey := "access_operations:access_requirement_code:" + code

	var ops []entity.AccessOperation
	err := a.cache.GetString(ctx, cacheKey, &ops)
	if err == nil {
		return ops, nil
	}

	if !errors.Is(err, goredis.Nil) {
		return nil, err
	}

	scopes := scope.GormScope().
		Join(fmt.Sprintf(
			"INNER JOIN %s AS aro ON aro.access_operation_id = %s.id",
			entity.AccessRequirementOperationTableName,
			entity.AccessOperationsTableName,
		)).
		Join(fmt.Sprintf(
			"INNER JOIN %s AS ar ON ar.id = aro.access_requirement_id",
			entity.AccessRequirementTableName,
		)).
		Where("ar.code = ?", code)

	ops, err = a.accessOpRepo.FindByConditions(ctx, nil, scopes)
	if err != nil {
		return nil, err
	}

	opsBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	err = a.cache.SetString(ctx, cacheKey, opsBytes, time.Hour)
	return ops, err
}
