package repository

import (
	"context"
	"github.com/dpe27/es-krake/internal/domain/product/entity"
)

const AttributeTypeTableName = "attribute_types"

type AttributeTypeRepository interface {
	TakeByID(ctx context.Context, ID int) (entity.AttributeType, error)
	GetAsDictionary(ctx context.Context) (map[int]string, error)
}
