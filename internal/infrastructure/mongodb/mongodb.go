package mdb

import (
	"context"
	"fmt"
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	mongolog "github.com/dpe27/es-krake/pkg/log/mongo"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	DB     *mongo.Client
	DBName string
	logger *log.Logger

	connAttempts int
	poolSize     int
	timeout      time.Duration
	connTimeout  time.Duration
	maxIdleTime  time.Duration
}

var (
	mongoInstance *Mongo
	once          sync.Once
)

func NewOrGetSingleton(cfg *config.Config) *Mongo {
	once.Do(func() {
		m, err := initMongo(cfg)
		if err != nil {
			panic(err)
		}
		mongoInstance = m
	})
	return mongoInstance
}

func initMongo(cfg *config.Config) (*Mongo, error) {
	m := &Mongo{
		logger:       log.With("service", "mongodb"),
		connAttempts: cfg.MDB.ConnAttempts,
		poolSize:     cfg.MDB.PoolSize,
		timeout:      time.Duration(cfg.MDB.Timeout) * time.Millisecond,
		connTimeout:  time.Duration(cfg.MDB.ConnTimeout) * time.Millisecond,
		maxIdleTime:  time.Duration(cfg.MDB.MaxIdleTime) * time.Millisecond,
		DBName:       cfg.MDB.Database,
	}

	mongoLogger := mongolog.NewMongoLog()
	loggerOptions := options.
		Logger().
		SetSink(mongoLogger).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	uri := fmt.Sprintf("mongodb://%s:%s", cfg.MDB.Hostname, cfg.MDB.Port)
	credential := options.Credential{
		AuthSource: cfg.MDB.AuthSource,
		Username:   cfg.MDB.Username,
		Password:   cfg.MDB.Password,
	}

	clientOpts := options.Client().
		ApplyURI(uri).
		SetAuth(credential).
		SetLoggerOptions(loggerOptions)

	for m.connAttempts > 0 {
		client, err := mongo.Connect(clientOpts)
		if err == nil {
			m.DB = client
			break
		}

		m.logger.Warn(
			context.Background(),
			"MongoDB is trying to connect",
			"error", err.Error(),
			"attempts left", m.connAttempts,
		)
		time.Sleep(3 * time.Millisecond)
		m.connAttempts--
	}

	if m.DB == nil {
		return nil, fmt.Errorf("MongoDB (initDB): connection failed")
	}

	return m, nil
}

func (m *Mongo) Ping(ctx context.Context) error {
	return m.DB.Ping(ctx, readpref.PrimaryPreferred())
}

func (m *Mongo) Close(ctx context.Context) {
	m.logger.Info(ctx, "Closing MongoDB")
	if err := m.DB.Disconnect(ctx); err != nil {
		m.logger.Error(ctx, "Error while closing MongoDB", "error", err.Error())
	}
}
