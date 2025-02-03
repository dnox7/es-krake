package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
)

const AttributeTypeTableName = "attribute_types"

type IAttributeTypeRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.AttributeType, error)
	GetAsDictionary(ctx context.Context) (map[int]string, error)
}
