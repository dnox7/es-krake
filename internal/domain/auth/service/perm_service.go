package service

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
)

type PermissionService interface {
	GetPermissionsWithRoleID(ctx context.Context, roleID int) ([]entity.Permission, error)
}
