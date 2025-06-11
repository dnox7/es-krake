package repository

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/product/entity"
	domainRepo "github.com/dpe27/es-krake/internal/domain/product/repository"
	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type productMetaRepo struct {
	logger *log.Logger
	db     *mdb.Mongo
}

func NewProductMetaRepository(db *mdb.Mongo) domainRepo.ProductMetaRespository {
	return &productMetaRepo{
		logger: log.With("repository", "product_meta_repo"),
		db:     db,
	}
}

// Create implements repository.ProductMetaRespository.
func (p *productMetaRepo) Create(
	ctx context.Context,
	attributes map[string]interface{},
) (entity.ProductMeta, error) {
	prodMeta := entity.ProductMeta{}
	err := utils.MapToStruct(attributes, &prodMeta)
	if err != nil {
		p.logger.Error(ctx, utils.ErrorMapToStruct, "error", err.Error())
		return entity.ProductMeta{}, err
	}

	coll := p.db.Cli().
		Database(p.db.DBName).
		Collection(entity.ProductMetaTableName)

	prodMeta.PrepareForInsert()
	res, err := coll.InsertOne(ctx, prodMeta)
	if err != nil {
		return entity.ProductMeta{}, err
	}

	prodMeta.ID = res.InsertedID.(primitive.ObjectID)
	return prodMeta, nil
}

// FindByConditions implements repository.ProductMetaRespository.
func (p *productMetaRepo) FindByConditions(
	ctx context.Context,
	filter interface{},
	spec specification.Base,
) ([]entity.ProductMeta, error) {
	opts, err := mdb.ToOptsBuilder[*options.FindOptionsBuilder](spec)
	if err != nil {
		p.logger.Error(ctx, err.Error())
		return nil, err
	}

	coll := p.db.Cli().
		Database(p.db.DBName).
		Collection(entity.ProductMetaTableName)

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
func (p *productMetaRepo) TakeByConditions(
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
	coll := p.db.Cli().
		Database(p.db.DBName).
		Collection(entity.ProductMetaTableName)

	err = coll.FindOne(ctx, filter, opts).Decode(&prodMeta)
	return prodMeta, err
}

// Update implements repository.ProductMetaRespository.
func (p *productMetaRepo) UpdateByID(
	ctx context.Context,
	ID interface{},
	operation interface{},
) (entity.ProductMeta, error) {
	idUpdate, ok := ID.(primitive.ObjectID)
	if !ok {
		return entity.ProductMeta{}, fmt.Errorf("Document's ID is invalid")
	}
	idFilter := bson.D{{Key: "_id", Value: idUpdate}}

	coll := p.db.Cli().
		Database(p.db.DBName).
		Collection(entity.ProductMetaTableName)
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updatedDoc := entity.ProductMeta{}
	err := coll.FindOneAndUpdate(ctx, idFilter, operation, opt).Decode(&updatedDoc)
	return updatedDoc, err
}
