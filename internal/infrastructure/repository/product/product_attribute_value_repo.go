package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	"github.com/Masterminds/squirrel"
)

type productAttributeValueRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

// func NewProductAttributeValueRepository(pg *db.PostgreSQL) domainRepo.IProductAttributeValueRepository {
// 	return &productAttributeValueRepository{
// 		logger: log.With("repo", "product_attribute_value_repo"),
// 		pg:     pg,
// 	}
// }

func (r *productAttributeValueRepository) TakeByID(ctx context.Context, ID int) (entity.ProductAttributeValue, error) {
	var prodAttrVal entity.ProductAttributeValue

	sql, args, err := r.pg.Builder.
		Select(
			"pav.id", "pav.product_id", "pav.attribute_id", "pav.value", "pav.created_at", "pav.updated_at",
			"attr.id", "attr.name",
			"aty.id", "aty.name",
		).
		From(domainRepo.ProductAttributeValueTableName).Suffix("AS pav").
		InnerJoin(domainRepo.AttributeTableName).Suffix("AS attr ON attr.id = pav.attribute_id").
		InnerJoin(domainRepo.AttributeTypeTableName).Suffix("AS aty ON aty.id = attr.attribute_type_id").
		Where(squirrel.Eq{"pav.id": ID}).
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return prodAttrVal, err
	}

	row := r.pg.DB.QueryRowxContext(ctx, sql, args...)
	if row.Err() != nil {
		return prodAttrVal, err
	}

}
