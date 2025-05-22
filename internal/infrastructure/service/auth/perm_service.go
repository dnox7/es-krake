package service

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	domainService "github.com/dpe27/es-krake/internal/domain/auth/service"
	scope "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/scope"
	"github.com/dpe27/es-krake/pkg/log"
)

type permService struct {
	permRepo repository.PermissionRepository
	logger   *log.Logger
}

func NewPermissionService(
	permRepo repository.PermissionRepository,
) domainService.PermissionService {
	return &permService{
		permRepo,
		log.With("service", "permission_serivce"),
	}
}

// GetPermissionsWithRoleID implements service.PermissionService.
func (p *permService) GetPermissionsWithRoleID(
	ctx context.Context,
	roleID int,
) ([]entity.Permission, error) {
	scopes := scope.GormScope().
		Join(fmt.Sprintf(
			"INNER JOIN %s AS rp ON rp.permission_id = %s.id",
			repository.RolePermissionTableName,
			repository.PermissionTableName,
		)).
		Join(fmt.Sprintf(
			"INNER JOIN %s AS r ON r.id = rp.role_id",
			repository.RoleTableName,
		)).
		Where("r.id = ?", roleID).
		Preload("Operations").
		Preload("Operations.AccessOperation")
	return p.permRepo.FindByConditions(ctx, nil, scopes)
}
