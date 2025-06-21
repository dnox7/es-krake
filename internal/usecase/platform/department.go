package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dpe27/es-krake/internal/domain/platform/cachekey"
	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	"github.com/redis/go-redis/v9"
)

func (pu *PlatformUsecase) GetDepartmentByID(ctx context.Context, ID int) (entity.Department, error) {
	key := cachekey.DepartmentByID(ID)

	var department entity.Department
	err := pu.deps.Cache.GetJSON(ctx, key, &department)
	if err == nil {
		return department, nil
	}

	if !errors.Is(err, redis.Nil) {
		return entity.Department{}, err
	}

	department, err = pu.deps.DepartmentRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": ID,
	}, nil)
	if err != nil {
		return entity.Department{}, err
	}

	_ = pu.deps.Cache.SetJSON(ctx, key, department, time.Hour)
	return department, nil
}
