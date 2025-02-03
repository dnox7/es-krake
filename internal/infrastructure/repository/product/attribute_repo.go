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

func (r *attributeRepository) CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.Attribute, error) {
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

	tx, err := r.pg.DB.Beginx()
	if err != nil {
		return entity.Attribute{}, err
	}

	committed := false
	defer (func() {
		if !committed {
			err := tx.Rollback()
			if err != nil {
				r.logger.Error("Cannot rollback transaction", "method", "attribute_repo - create_with_tx", "detail", err)
			}
		}
	})()

	err = tx.QueryRowx(sql, args...).StructScan(&attributeEntity)
	if err != nil {
		return entity.Attribute{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.Attribute{}, err
	}

	return attributeEntity, nil
}
