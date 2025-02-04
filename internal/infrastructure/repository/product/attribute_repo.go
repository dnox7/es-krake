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

type attributeRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewAttributeRepository(pg *db.PostgreSQL) domainRepo.IAttributeRepository {
	return &attributeRepository{
		logger: log.With("repo", "attribute_repo"),
		pg:     pg,
	}
}

func (r *attributeRepository) TakeByID(ctx context.Context, ID int) (entity.Attribute, error) {
	var attribute entity.Attribute

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From(domainRepo.AttributeTableName).
		Where(sq.Eq{"id": ID}).
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg, "detail", err)
		return attribute, err
	}

	err = r.pg.DB.GetContext(ctx, &attribute, sql, args...)
	return attribute, err
}

func (r *attributeRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Attribute, error) {
	var attribute entity.Attribute

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From(domainRepo.AttributeTableName).
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg, "detail", err)
		return attribute, err
	}

	err = r.pg.DB.GetContext(ctx, &attribute, sql, args...)
	return attribute, err
}

func (r *attributeRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Attribute, error) {
	attributes := []entity.Attribute{}

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From(domainRepo.AttributeTableName).
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg, "detail", err)
		return attributes, err
	}

	err = r.pg.DB.SelectContext(ctx, &attributes, sql, args...)
	return attributes, err
}

func (r *attributeRepository) Create(ctx context.Context, attributes map[string]interface{}) (entity.Attribute, error) {
	var attributeEntity entity.Attribute

	if err := utils.MapToStruct(attributes, &attributeEntity); err != nil {
		return attributeEntity, err
	}

	sql, args, err := r.pg.Builder.
		Insert(domainRepo.AttributeTableName).
		Columns("name", "description", "attribute_type_id", "is_required", "display_order").
		Values(attributeEntity.Name, attributeEntity.Description, attributeEntity.AttributeTypeID, attributeEntity.IsRequired, attributeEntity.DisplayOrder).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg, "detail", err)
		return attributeEntity, err
	}

	err = utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&attributeEntity)
	})

	if err != nil {
		return entity.Attribute{}, err
	}

	return attributeEntity, nil
}

func (r *attributeRepository) Update(
	ctx context.Context,
	attributeEntity entity.Attribute,
	attributesToUpdate map[string]interface{},
) (entity.Attribute, error) {
	if err := utils.MapToStruct(attributesToUpdate, &attributeEntity); err != nil {
		return attributeEntity, err
	}

	sql, args, err := r.pg.Builder.
		Update(domainRepo.AttributeTableName).
		SetMap(map[string]interface{}{
			"name":              attributeEntity.Name,
			"description":       attributeEntity.Description,
			"attribute_type_id": attributeEntity.AttributeTypeID,
			"is_required":       attributeEntity.IsRequired,
			"display_order":     attributeEntity.DisplayOrder,
		}).
		Where(sq.Eq{"id": attributeEntity.ID}).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		r.logger.ErrorContext(ctx, utils.ErrQueryBuilderFailedMsg, "detail", err)
		return entity.Attribute{}, err
	}

	err = utils.Transaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&attributeEntity)
	})

	if err != nil {
		return entity.Attribute{}, err
	}

	return attributeEntity, nil
}
