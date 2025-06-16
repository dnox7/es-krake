package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/cachekey"
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
	goredis "github.com/redis/go-redis/v9"
)

type accessOperationService struct {
	accessOpRepo repository.AccessOperationRepository
	redisRepo    redis.RedisRepository
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
	key := cachekey.OperationsByAccessRequirementCode(code)

	var ops []entity.AccessOperation
	err := a.redisRepo.GetJSON(ctx, key, &ops)
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

	_ = a.redisRepo.SetJSON(ctx, key, ops, time.Hour)
	return ops, nil
}

// GetOperationsByPermissionID implements service.AccessOperationService.
func (a *accessOperationService) GetOperationsByPermissionID(
	ctx context.Context,
	permissionID int,
) ([]entity.AccessOperation, error) {
	key := cachekey.OperationsByPermissionID(permissionID)

	var ops []entity.AccessOperation
	err := a.redisRepo.GetJSON(ctx, key, &ops)
	if err == nil {
		return ops, nil
	}

	if !errors.Is(err, goredis.Nil) {
		return nil, err
	}

	scopes := scope.GormScope().
		Join(fmt.Sprintf(
			"INNER JOIN %s AS aro ON aro.access_operation_id = %s.id",
			entity.PermissionOperationTableName,
			entity.AccessOperationsTableName,
		)).
		Where("permission_id = ?", permissionID)

	ops, err = a.accessOpRepo.FindByConditions(ctx, nil, scopes)
	if err != nil {
		return nil, err
	}

	_ = a.redisRepo.SetJSON(ctx, key, ops, time.Hour)
	return ops, nil
}
