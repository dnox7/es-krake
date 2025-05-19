package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type RoleType string

const (
	RoleTypeTableName = "role_types"

	PlatformRoleType   RoleType = "plarform"
	EnterpriseRoleType RoleType = "enterprise"
)

type RoleTypeRepository interface {
	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.RoleType, error)
}
