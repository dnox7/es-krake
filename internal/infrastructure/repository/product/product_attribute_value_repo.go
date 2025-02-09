package repository

import (
	"context"
	"database/sql"
	"pech/es-krake/internal/domain"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type productAttributeValueRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewProductAttributeValueRepository(pg *db.PostgreSQL) domainRepo.IProductAttributeValueRepository {
	return &productAttributeValueRepository{
		logger: log.With("repo", "product_attribute_value_repo"),
		pg:     pg,
	}
}

func (r *productAttributeValueRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.ProductAttributeValue, error) {
	var prodAttrVal entity.ProductAttributeValue

	query, args, err := r.getProdAttrValQueryBuilder(conditions)
	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return prodAttrVal, err
	}

	row := r.pg.DB.QueryRowxContext(ctx, query, args)
	err = row.Err()
	if err == sql.ErrNoRows {
		err = domain.ErrRecordNotFound
	}

	if err != nil {
		return prodAttrVal, err
	}

	prodAttrVal.Attribute = &entity.Attribute{}
	prodAttrVal.Type = &entity.AttributeType{}

	err = row.Scan(
		&prodAttrVal.ID,
		&prodAttrVal.ProductID,
		&prodAttrVal.AttributeID,
		&prodAttrVal.Value,
		&prodAttrVal.CreatedAt,
		&prodAttrVal.UpdatedAt,
		&prodAttrVal.Attribute.ID,
		&prodAttrVal.Attribute.Name,
		&prodAttrVal.Type.ID,
		&prodAttrVal.Type.Name,
	)
	return prodAttrVal, err

}

func (r *productAttributeValueRepository) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entity.ProductAttributeValue, error) {
	var pavList []entity.ProductAttributeValue

	query, args, err := r.getProdAttrValQueryBuilder(conditions)
	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return nil, err
	}

	rows, err := r.pg.DB.QueryxContext(ctx, query, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pav := entity.ProductAttributeValue{
			Attribute: &entity.Attribute{},
			Type:      &entity.AttributeType{},
		}

		if err := rows.Scan(
			&pav.ID,
			&pav.ProductID,
			&pav.AttributeID,
			&pav.Value,
			&pav.CreatedAt,
			&pav.UpdatedAt,
			&pav.Attribute.ID,
			&pav.Attribute.Name,
			&pav.Type.ID,
			&pav.Type.Name,
		); err != nil {
			return nil, err
		}

		pavList = append(pavList, pav)
	}
	return pavList, nil
}

func (r *productAttributeValueRepository) getProdAttrValQueryBuilder(
	conditions map[string]interface{},
) (string, interface{}, error) {
	return r.pg.Builder.
		Select(
			"pav.id",
			"pav.product_id",
			"pav.attribute_id",
			"pav.value",
			"pav.created_at",
			"pav.updated_at",
			"attr.id",
			"attr.name",
			"aty.id",
			"aty.name",
		).
		From(domainRepo.ProductAttributeValueTableName + " AS pav").
		InnerJoin(domainRepo.AttributeTableName + " AS attr ON attr.id = pav.attribute_id").
		InnerJoin(domainRepo.AttributeTypeTableName + " AS aty ON aty.id = attr.attribute_type_id").
		Where(sq.Eq(conditions)).
		ToSql()
}

func (r *productAttributeValueRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.ProductAttributeValue, error) {
	var pav entity.ProductAttributeValue

	if err := utils.MapToStruct(attributes, &pav); err != nil {
		return pav, err
	}

	sql, args, err := r.pg.Builder.
		Insert(domainRepo.ProductAttributeValueTableName).
		Columns("product_id", "attribute_id", "value").
		Values(pav.ProductID, pav.AttributeID, pav.Value).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return pav, err
	}

	err = utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&pav)
	})

	return pav, err
}

func (r *productAttributeValueRepository) Update(
	ctx context.Context,
	pav entity.ProductAttributeValue,
	attributesToUpdate map[string]interface{},
) (entity.ProductAttributeValue, error) {
	if err := utils.MapToStruct(attributesToUpdate, &pav); err != nil {
		return pav, err
	}

	sql, args, err := r.pg.Builder.
		Update(domainRepo.ProductAttributeValueTableName).
		SetMap(map[string]interface{}{
			"product_id":   pav.ProductID,
			"attribute_id": pav.AttributeID,
			"value":        pav.Value,
		}).
		Where(sq.Eq{"id": pav.ID}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return pav, err
	}

	err = utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&pav)
	})

	return pav, err
}

func (r *productAttributeValueRepository) DeleteByConditions(ctx context.Context, conditions map[string]interface{}) error {
	sql, args, err := r.pg.Builder.
		Delete(domainRepo.ProductAttributeValueTableName).
		Where(squirrel.Eq(conditions)).
		ToSql()
	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg)
		return err
	}

	return utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, sql, args...)
		r.logger.WarnContext(ctx, "Deleting productAttributeValue")
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		r.logger.WarnContext(ctx, "Deleted productAttributeValue", rowsAffected)
		return err
	})
}
