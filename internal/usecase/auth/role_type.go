package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/cachekey"
	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/redis/go-redis/v9"
)

func (au *AuthUsecase) GetAllRoleTypes(ctx context.Context) ([]entity.RoleType, error) {
	key := cachekey.AllRoleTypesKey()

	var roleTypes []entity.RoleType
	err := au.deps.Cache.GetJSON(ctx, key, &roleTypes)
	if err == nil {
		return roleTypes, nil
	}

	if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	roleTypes, err = au.deps.RoleTypeRepo.FindByConditions(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	_ = au.deps.Cache.SetJSON(ctx, key, roleTypes, time.Hour)
	return roleTypes, nil
}

func (au *AuthUsecase) GetRoleTypeByID(ctx context.Context, ID int) (entity.RoleType, error) {
	key := cachekey.RoleTypeByIDKey(ID)

	var roleType entity.RoleType
	err := au.deps.Cache.GetJSON(ctx, key, &roleType)
	if err == nil {
		return roleType, nil
	}

	if !errors.Is(err, redis.Nil) {
		return entity.RoleType{}, err
	}

	roleType, err = au.deps.RoleTypeRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": ID,
	}, nil)
	if err != nil {
		return entity.RoleType{}, err
	}

	_ = au.deps.Cache.SetJSON(ctx, key, roleType, time.Hour)
	return roleType, nil
}
