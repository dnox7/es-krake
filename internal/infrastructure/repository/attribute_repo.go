package repository

import (
	"context"
	"fmt"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"

	sq "github.com/Masterminds/squirrel"
)

type attributeRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewAttributeRepository(l *log.Logger, pg *db.PostgreSQL) *domainRepo.IAttributeRepository {
	return &attributeRepository{
		logger: l,
		pg:     pg,
	}
}

func (r *attributeRepository) TakeByID(ctx context.Context, ID int) (entity.Attribute, error) {
	var attribute entity.Attribute

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From(attribute.TableName()).
		Where("id", ID).
		ToSql()

	if err != nil {
		return attribute, fmt.Errorf("query build failed")
	}

	err = r.pg.DB.GetContext(ctx, &attribute, sql, args...)
	return attribute, err
}

func (r *attributeRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}) (entity.Attribute, error) {
	var attribute entity.Attribute

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From(attribute.TableName()).
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		return attribute, fmt.Errorf("query build failed")
	}

	err = r.pg.DB.GetContext(ctx, &attribute, sql, args...)
	return attribute, err
}

func (r *attributeRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}) ([]entity.Attribute, error) {
	attributes := []entity.Attribute{}

	sql, args, err := r.pg.Builder.
		Select("id", "name", "description", "is_required", "display_order", "created_at", "updated_at").
		From("attributes").
		Where(sq.Eq(conditions)).
		ToSql()

	if err != nil {
		r.logger.Error("Query build failed", "detail", err)
		return attributes, nil
	}

	err = r.pg.DB.SelectContext(ctx, &attributes, sql, args...)
	return attributes, err
}
