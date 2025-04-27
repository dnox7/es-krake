package mdb

import (
	"context"
	"pech/es-krake/internal/domain/shared/transaction"
	"pech/es-krake/pkg/log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoTx struct {
	tx *mongo.Client
}

func NewMongoTx(tx *mongo.Client) transaction.Base {
	return &mongoTx{tx}
}

func (m *mongoTx) GetTx() interface{} {
	return m.tx
}

func MongoTransaction(
	ctx context.Context,
	l *log.Logger,
	cli *mongo.Client,
	opts *options.TransactionOptionsBuilder,
	fn func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	session, err := cli.StartSession()
	if err != nil {
		l.Error(ctx, "Failed to begin a Mongo transaction")
		return nil, err
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, fn, opts)
	if err != nil {
		return nil, err
	}

	return res, nil
}
