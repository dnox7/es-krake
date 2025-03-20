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

type productRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewProductRepository(pg *db.PostgreSQL) domainRepo.IProductRepository {
	return &productRepository{
		logger: log.With("repo", "product_repo"),
		pg:     pg,
	}
}

func (r *productRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Product, error) {
	var prod entity.Product

	sql, args, err := r.pg.Builder.
		Select(
			"id",
			"name",
			"sku",
			"description",
			"price",
			"has_options",
			"is_allowed_to_order",
			"is_publised",
			"is_featured",
			"is_visible_individually",
			"stock_tracking_enabled",
			"stock_quantity",
			"tax_class_id",
			"meta_title",
			"meta_keyword",
		).
		From(domainRepo.ProductTableName).
		Where(sq.Eq(conditions)).
		ToSql()
	if err != nil {
		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
		return prod, err
	}

	err = r.pg.DB.GetContext(ctx, &prod, sql, args...)
	return prod, err
}

func (r *productRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Product, error) {
	return nil, nil
}

func (r *productRepository) Create(ctx context.Context, attributes map[string]interface{}) (entity.Product, error) {
	return entity.Product{}, nil
}

func (r *productRepository) CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Product, error) {
	panic("lmao")
}

func (r *productRepository) UpdateWithTx(ctx context.Context, prod entity.Product, attributesToUpdate map[string]interface{}) (entity.Product, error) {
	panic("lmao")
}
