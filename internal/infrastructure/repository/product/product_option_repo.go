package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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

func (r *productOptionRepository) TakeByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) (entity.ProductOption, error) {
	var productOption entity.ProductOption

	sql, args, err := r.pg.Builder.
		Select("id", "product_id", "name", "description", "created_at", "updated_at").
		From(domainRepo.ProductOptionTableName).
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.Error(utils.ErrQueryBuilderFailedMsg)
		return productOption, err
	}

	err = r.pg.DB.GetContext(ctx, &productOption, sql, args...)

	return productOption, err
}

func (r *productOptionRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entity.ProductOption, error) {
	sql, args, err := r.pg.Builder.
		Select(
			"po.id", "po.product_id", "po.name", "po.description", "po.created_at", "po.updated_at",
			"pav.id", "pav.value",
			"pa.id", "pav.name",
		).
		From(domainRepo.ProductOptionTableName + " AS po").
		InnerJoin(domainRepo.ProductAttributeValueTableName + " AS pav ON pav.product_option_id = po.id").
		InnerJoin(domainRepo.AttributeTableName + " AS pa ON pa.id = pav.attribute_id").
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.Error(utils.ErrQueryBuilderFailedMsg)
		return nil, err
	}

	rows, err := r.pg.DB.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	optionsMap := make(map[int]*entity.ProductOption)

	for rows.Next() {
		opt := entity.ProductOption{}
		attribute := entity.Attribute{}
		attributeValue := entity.ProductAttributeValue{}

		if err := rows.Scan(
			&opt.ID, &opt.ProductID, &opt.Name, &opt.Description, &opt.CreatedAt, &opt.UpdatedAt,
			&attributeValue.ID, &attributeValue.Value,
			&attribute.ID, &attribute.Name,
		); err != nil {
			return nil, err
		}

		attributeValue.Attribute = &attribute
		if _, ok := optionsMap[opt.ID]; !ok {
			opt.Attributes = []*entity.ProductAttributeValue{&attributeValue}
			optionsMap[opt.ID] = &opt
		} else {
			optionsMap[opt.ID].Attributes = append(optionsMap[opt.ID].Attributes, &attributeValue)
		}
	}

	optionsSlice := []entity.ProductOption{}
	for _, opt := range optionsMap {
		optionsSlice = append(optionsSlice, *opt)
	}

	return optionsSlice, nil
}

func (r *productOptionRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.ProductOption, error) {
	var opt entity.ProductOption

	if err := utils.MapToStruct(attributes, &opt); err != nil {
		return opt, err
	}

	sql, args, err := r.pg.Builder.
		Insert(domainRepo.ProductOptionTableName).
		Columns("name", "description", "product_id").
		Values(opt.Name, opt.Description, opt.ProductID).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		r.logger.Error(utils.ErrQueryBuilderFailedMsg)
		return opt, err
	}

	err = r.pg.DB.QueryRowxContext(ctx, sql, args...).Scan(&opt)
	return opt, err
}

func (r *productOptionRepository) UpdateWithTx(
	ctx context.Context,
	option entity.ProductOption,
	attributesToUpdate map[string]interface{},
) (entity.ProductOption, error) {
	if err := utils.MapToStruct(attributesToUpdate, &option); err != nil {
		return option, err
	}

	sql, args, err := r.pg.Builder.
		Update(domainRepo.ProductOptionTableName).
		SetMap(map[string]interface{}{
			"name":        option.Name,
			"description": option.Description,
			"product_id":  option.ProductID,
		}).
		Where(sq.Eq{"id": option.ID}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		r.logger.Error(utils.ErrQueryBuilderFailedMsg)
		return option, err
	}

	err = utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		return tx.QueryRowxContext(ctx, sql, args...).Scan(&option)
	})

	return option, err
}
