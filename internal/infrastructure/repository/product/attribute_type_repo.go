package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/rdb"
	"pech/es-krake/pkg/log"
)

type attributeTypeRepo struct {
	logger *log.Logger
	pg     *rdb.PostgreSQL
}

func NewAttributeTypeRepository(pg *rdb.PostgreSQL) domainRepo.AttributeTypeRepository {
	return &attributeTypeRepo{
		logger: log.With("repository", "attribute_type_repo"),
		pg:     pg,
	}
}

// TakeByID implements repository.AttributeTypeRepository.
func (r *attributeTypeRepo) TakeByID(ctx context.Context, ID int) (entity.AttributeType, error) {
	attrType := entity.AttributeType{}
	db := r.pg.DB.WithContext(ctx)
	err := db.Take(&attrType, ID).Error
	return attrType, err
}

// GetAsDictionary implements repository.AttributeTypeRepository.
func (r *attributeTypeRepo) GetAsDictionary(ctx context.Context) (map[int]string, error) {
	var attrTypes []entity.AttributeType

	err := r.pg.DB.
		WithContext(ctx).
		Table(domainRepo.AttributeTypeTableName).
		Select("id", "name").
		Find(&attrTypes).Error
	if err != nil {
		return nil, err
	}

	dict := make(map[int]string)
	for _, attrType := range attrTypes {
		dict[attrType.ID] = attrType.Name
	}
	return dict, nil
}
