package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

const EnterpriseAccountTableName = "enterprise_accounts"

type EnterpriseAccountRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.EnterpriseAccount, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.EnterpriseAccount, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.EnterpriseAccount, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.EnterpriseAccount, error)

	Update(ctx context.Context, acc entity.EnterpriseAccount, attributesToUpdate map[string]interface{}) (entity.EnterpriseAccount, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, acc entity.EnterpriseAccount, attributesToUpdate map[string]interface{}) (entity.EnterpriseAccount, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
