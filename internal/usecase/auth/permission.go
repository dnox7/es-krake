package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/cachekey"
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/redis/go-redis/v9"
)

func (au *AuthUsecase) GetPermissionByID(ctx context.Context, ID int) (entity.Permission, error) {
	key := cachekey.PermByIDKey(ID)
	var perm entity.Permission
	err := au.deps.Cache.GetJSON(ctx, key, &perm)
	if err == nil {
		return perm, nil
	}

	if !errors.Is(err, redis.Nil) {
		return entity.Permission{}, err
	}

	perm, err = au.deps.PermissionRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": ID,
	}, nil)
	if err != nil {
		return entity.Permission{}, err
	}

	_ = au.deps.Cache.SetJSON(ctx, key, perm, time.Hour)
	return perm, nil
}

func (au *AuthUsecase) GetPermissionsWithRoleID(ctx context.Context, roleID int) ([]entity.Permission, error) {
	return au.deps.PermissionService.GetPermissionsWithRoleID(ctx, roleID)
}
