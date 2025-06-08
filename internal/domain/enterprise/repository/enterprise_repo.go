package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const EnterpriseTableName = "enterprises"

type EnterpriseRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Enterprise, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Enterprise, error)

	PluckIDByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]int, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Enterprise, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Enterprise, error)

	Update(ctx context.Context, ent entity.Enterprise, attributesToUpdate map[string]interface{}) (entity.Enterprise, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, ent entity.Enterprise, attributesToUpdate map[string]interface{}) (entity.Enterprise, error)
}
