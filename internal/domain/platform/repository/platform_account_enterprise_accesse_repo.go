package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/platform/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

type PlatformAccountEnterpriseAccessRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.PlatformAccountEnterpriseAccess, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.PlatformAccountEnterpriseAccess, error)

	Update(ctx context.Context, pae entity.PlatformAccountEnterpriseAccess, attributesToUpdate map[string]interface{}) (entity.PlatformAccountEnterpriseAccess, error)

	DeleteByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) error
}
