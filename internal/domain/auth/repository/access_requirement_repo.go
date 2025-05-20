package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const AccessRequirementTableName = "access_requirements"

type AccessRequirementRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.AccessRequirement, error)
}
