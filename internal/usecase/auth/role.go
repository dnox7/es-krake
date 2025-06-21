package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/cachekey"
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/redis/go-redis/v9"
)

func (au *AuthUsecase) GetRoleByIDWithPermissions(ctx context.Context, ID int) (entity.Role, error) {
	key := cachekey.RoleByIDKey(ID)

	var role entity.Role
	err := au.deps.Cache.GetJSON(ctx, key, &role)
	if err == nil {
		return role, nil
	}

	if !errors.Is(err, redis.Nil) {
		return entity.Role{}, err
	}

	role, err = au.deps.RoleRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": ID,
	}, scope.GormScope().Preload("RolePermissions.Permission"))
	if err != nil {
		return entity.Role{}, err
	}

	_ = au.deps.Cache.SetJSON(ctx, key, role, time.Hour)
	return role, nil
}
