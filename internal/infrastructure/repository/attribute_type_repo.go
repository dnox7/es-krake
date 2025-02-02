package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
)

type attributeTypeRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewAttributeTypeRepository(l *log.Logger, pg *db.PostgreSQL) domainRepo.IAttributeTypeRepository {
	return &attributeTypeRepository{
		logger: l,
		pg:     pg,
	}
}

func (r *attributeTypeRepository) TakeByID(ctx context.Context, ID int) (entity.AttributeType, error) {
	var attributeType entity.AttributeType

	sql, args, err := r.pg.Builder.
		Select("id, name, created_at, updated_at").
		From(attributeType.TableName()).
		Where("id", ID).
		ToSql()

	if err != nil {
		r.logger.Error("Query build failed", "detail", err)
		return attributeType, err
	}

	err = r.pg.DB.GetContext(ctx, &attributeType, sql, args...)
	return attributeType, err
}

func (r *attributeTypeRepository) GetAsDictionary(ctx context.Context) (map[int]string, error) {
	var attributeTypes []entity.AttributeType

	sql, args, err := r.pg.Builder.
		Select("id", "name").
		From("attributes").
		ToSql()

	if err != nil {
		r.logger.Error("Query build failed", "detail", err)
		return nil, err
	}

	err = r.pg.DB.SelectContext(ctx, &attributeTypes, sql, args...)
	if err != nil {
		return nil, err
	}

	dict := make(map[int]string)
	for _, attributeType := range attributeTypes {
		dict[attributeType.ID] = attributeType.Name
	}
	return dict, nil
}
