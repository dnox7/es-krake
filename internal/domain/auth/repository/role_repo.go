package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const RoleTableName = "roles"

type RoleRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Role, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Role, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Role, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, role entity.Role, attributesToUpdate map[string]interface{}) (entity.Role, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
