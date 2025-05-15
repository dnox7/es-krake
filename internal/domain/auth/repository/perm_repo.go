package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const PermissionTableName = "permissions"

type PermissionRepository interface {
	TakeByCondition(ctx context.Context, condition map[string]interface{}, spec specification.Base) (entity.Permission, error)
	FindByCondition(ctx context.Context, condition map[string]interface{}, spec specification.Base) ([]entity.Permission, error)
	Create(ctx context.Context, attributes map[string]interface{}) (entity.Permission, error)
	Update(ctx context.Context, perm entity.Permission, attributesToUpdate map[string]interface{}) (entity.Permission, error)
}
