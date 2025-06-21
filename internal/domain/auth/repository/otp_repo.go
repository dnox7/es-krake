package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/internal/domain/shared/transaction"
)

type OtpRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Otp, error)

	CreateWithTx(ctx context.Context, tx transaction.Base, attributes map[string]interface{}) (entity.Otp, error)

	UpdateWithTx(ctx context.Context, tx transaction.Base, otp entity.Otp, attributesToUpdate map[string]interface{}) (entity.Otp, error)

	DeleteByConditionsWithTx(ctx context.Context, tx transaction.Base, conditions map[string]interface{}, spec specification.Base) error
}
