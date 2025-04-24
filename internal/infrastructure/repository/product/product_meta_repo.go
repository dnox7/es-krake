package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/domain/shared/specification"
	mdb "pech/es-krake/internal/infrastructure/mongodb"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type productMetaRepository struct {
	logger *log.Logger
	db     *mdb.Mongo
}

func NewProductMetaRepository(db *mdb.Mongo) domainRepo.ProductMetaRespository {
	return &productMetaRepository{
		logger: log.With("repo", "product_meta_repo"),
		db:     db,
	}
}

// Create implements repository.ProductMetaRespository.
func (p *productMetaRepository) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.ProductMeta, error) {
	prodMeta := entity.ProductMeta{}
	err := utils.MapToStruct(attributes, &prodMeta)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.ProductMeta{}, err
	}

	coll := p.db.DB.
		Database(p.db.DBName).
		Collection(domainRepo.ProductMetaTableName)

	prodMeta.PrepareForInsert()
	res, err := coll.InsertOne(ctx, prodMeta)
	if err != nil {
		return entity.ProductMeta{}, err
	}

	prodMeta.ID = res.InsertedID.(primitive.ObjectID)
	return prodMeta, nil
}

// FindByConditions implements repository.ProductMetaRespository.
func (p *productMetaRepository) FindByConditions(
	ctx context.Context,
	filter interface{},
	spec specification.Base,
) ([]entity.ProductMeta, error) {
	opts, err := mdb.ToOptsBuilder[*options.FindOptionsBuilder](spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	coll := p.db.DB.
		Database(p.db.DBName).
		Collection(domainRepo.ProductMetaTableName)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		p.logger.Error(ctx, mdb.ErrGetMongoCursor, "error", err.Error())
		return nil, err
	}

	prodMetas := []entity.ProductMeta{}
	err = cursor.All(ctx, &prodMetas)
	if err != nil {
		return nil, err
	}
	return prodMetas, err
}

// TakeByConditions implements repository.ProductMetaRespository.
func (p *productMetaRepository) TakeByConditions(
	ctx context.Context,
	filter interface{},
	spec specification.Base,
) (entity.ProductMeta, error) {
	opts, err := mdb.ToOptsBuilder[*options.FindOneOptionsBuilder](spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return entity.ProductMeta{}, err
	}

	prodMeta := entity.ProductMeta{}
	coll := p.db.DB.
		Database(p.db.DBName).
		Collection(domainRepo.ProductMetaTableName)

	err = coll.FindOne(ctx, filter, opts).Decode(&prodMeta)
	return prodMeta, err
}
