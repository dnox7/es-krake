package service

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
)

type PermissionService interface {
	GetPermissionsWithKcUserID(ctx context.Context, kcUserID int) ([]entity.Permission, error)
}
