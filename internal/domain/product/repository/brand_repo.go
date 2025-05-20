package repository

import (
	"context"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
)

const BrandTableName = "brands"

type BrandRepository interface {
	TakeByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) (entity.Brand, error)

	FindByConditions(ctx context.Context, conditions map[string]interface{}, spec specification.Base) ([]entity.Brand, error)

	Create(ctx context.Context, attributes map[string]interface{}) (entity.Brand, error)

	Update(ctx context.Context, brand entity.Brand, attributesToUpdate map[string]interface{}) (entity.Brand, error)
}
