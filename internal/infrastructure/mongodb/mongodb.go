package mdb

import (
	"context"
	"fmt"
	"pech/es-krake/config"
	"pech/es-krake/pkg/log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	DB *mongo.Client

	connAttempts int
	poolSize     int
	timeout      time.Duration
	connTimeout  time.Duration
	maxIdleTime  time.Duration
}

func initMongo(cfg *config.Config) (*Mongo, error) {
	logger := log.With("service", "mongodb")

	m := &Mongo{
		connAttempts: cfg.RDB.ConnAttempts,
		poolSize:     cfg.MDB.PoolSize,
		timeout:      time.Duration(cfg.MDB.Timeout) * time.Millisecond,
		connTimeout:  time.Duration(cfg.MDB.ConnTimeout) * time.Millisecond,
		maxIdleTime:  time.Duration(cfg.MDB.MaxIdleTime) * time.Millisecond,
	}

	uri := fmt.Sprintf("mongodb://%s:%s", cfg.MDB.Hostname, cfg.RDB.Port)
	credential := options.Credential{
		AuthSource: cfg.MDB.AuthSource,
		Username:   cfg.MDB.Username,
		Password:   cfg.MDB.Password,
	}
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)

	for m.connAttempts > 0 {
		client, err := mongo.Connect(clientOpts)
		if err == nil {
			m.DB = client
			break
		}

		logger.Warn(
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
