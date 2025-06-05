package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type RoleType int

const (
	RoleTypeTableName = "role_types"

	PlatformRoleType   RoleType = 1
	EnterpriseRoleType RoleType = 2
)

type RoleTypeRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.RoleType, error)
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.RoleType, error)
}
