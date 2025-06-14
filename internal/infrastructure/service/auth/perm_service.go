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
	scope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/pkg/log"
	goredis "github.com/redis/go-redis/v9"
)

type permService struct {
	permRepo  repository.PermissionRepository
	redisRepo redis.RedisRepository
	logger    *log.Logger
}

func NewPermissionService(
	permRepo repository.PermissionRepository,
	redisRepo redis.RedisRepository,
) domainService.PermissionService {
	return &permService{
		permRepo,
		redisRepo,
		log.With("service", "permission_serivce"),
	}
}

// GetPermissionsWithRoleID implements service.PermissionService.
func (p *permService) GetPermissionsWithRoleID(
	ctx context.Context,
	roleID int,
) ([]entity.Permission, error) {
	key := cachekey.PermsByRoleID(roleID)

	var perms []entity.Permission
	err := p.redisRepo.GetJSON(ctx, key, &perms)
	if err == nil {
		return perms, nil
	}

	if !errors.Is(err, goredis.Nil) {
		return nil, err
	}

	scopes := scope.GormScope().
		Join(fmt.Sprintf(
			"INNER JOIN %s AS rp ON rp.permission_id = %s.id",
			entity.RolePermissionTableName,
			entity.PermissionTableName,
		)).
		Join(fmt.Sprintf(
			"INNER JOIN %s AS r ON r.id = rp.role_id",
			entity.RoleTableName,
		)).
		Where("r.id = ?", roleID).
		Preload("Operations").
		Preload("Operations.AccessOperation")

	perms, err = p.permRepo.FindByConditions(ctx, nil, scopes)
	if err != nil {
		return nil, err
	}

	_ = p.redisRepo.SetJSON(ctx, key, perms, time.Hour)
	return perms, nil
}
