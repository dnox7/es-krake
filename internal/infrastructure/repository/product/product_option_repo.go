package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	sq "github.com/Masterminds/squirrel"
)

type productOptionRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewProductOptionRepository(pg *db.PostgreSQL) domainRepo.IProductOptionRepository {
	return &productOptionRepository{
		logger: log.With("repo", "product_option_repo"),
		pg:     pg,
	}
}

func (r *productOptionRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductOption, error) {
	var productOption entity.ProductOption

	sql, args, err := r.pg.Builder.
		Select("po.id", "po.product_id", "po.name", "po.description", "po.created_at", "po.updated_at", ).
		From(domainRepo.ProductOptionTableName + "AS po").
        InnerJoin(domainRepo.ProductAttributeValueTableName + "AS pav ON pav.product_option_id = po.id").
        InnerJoin(domainRepo.AttributeTableName + "AS pa ON pav.attribute_id = pa.id")
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.Error(utils.ErrQueryBuilderFailedMsg)
		return productOption, err
	}

	err = r.pg.DB.GetContext(ctx, &productOption, sql, args...)

	return productOption, err
}

func (r *productOptionRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.ProductOption, error) {
    var 
}
