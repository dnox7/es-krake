package mdb

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	DBName string
	mu     sync.RWMutex
	cli    *mongo.Client
	logger *log.Logger
	params mongoParams
}

type mongoParams struct {
	uri          string
	authSource   string
	connAttempts int
	poolSize     int
	timeout      time.Duration
	connTimeout  time.Duration
	maxIdleTime  time.Duration

	loggerOpt *options.LoggerOptions
}

var (
	mongoInstance *Mongo
	once          sync.Once
)

func NewOrGetSingleton(ctx context.Context, cfg *config.Config, cred *config.MongoCredentials) *Mongo {
	once.Do(func() {
		m, err := initMongo(ctx, cfg, cred)
		if err != nil {
			panic(err)
		}
		mongoInstance = m
	})
	return mongoInstance
}

func initMongo(ctx context.Context, cfg *config.Config, cred *config.MongoCredentials) (*Mongo, error) {
	mongoLogger := newMongoLog()

	m := &Mongo{
		logger: log.With("service", "mongodb"),
		DBName: cfg.MDB.Database,
		params: mongoParams{
			uri:          "mongodb://" + net.JoinHostPort(cfg.MDB.Hostname, cfg.MDB.Port),
			authSource:   cfg.MDB.AuthSource,
			connAttempts: cfg.MDB.ConnAttempts,
			poolSize:     cfg.MDB.PoolSize,
			timeout:      time.Duration(cfg.MDB.Timeout) * time.Millisecond,
			connTimeout:  time.Duration(cfg.MDB.ConnTimeout) * time.Millisecond,
			maxIdleTime:  time.Duration(cfg.MDB.MaxIdleTime) * time.Millisecond,
		},
	}

	m.params.loggerOpt = options.
		Logger().
		SetSink(mongoLogger).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	if err := m.RetryConn(ctx, cred); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Mongo) setCli(newCli *mongo.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cli = newCli
}

func (m *Mongo) Cli() *mongo.Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cli
}

func (m *Mongo) RetryConn(ctx context.Context, cred *config.MongoCredentials) error {
	cliCred := options.Credential{
		AuthSource: m.params.authSource,
		Username:   cred.Username,
		Password:   cred.Password,
	}

	clientOpts := options.Client().
		ApplyURI(m.params.uri).
		SetAuth(cliCred).
		SetMaxPoolSize(uint64(m.params.poolSize)).
		SetTimeout(m.params.timeout).
		SetConnectTimeout(m.params.connTimeout).
		SetLoggerOptions(m.params.loggerOpt)

	prev := m.Cli()
	connAttempts := m.params.connAttempts
	for connAttempts > 0 {
		err := m.conn(clientOpts)
		if err == nil {
			break
		}
		m.logger.Warn(ctx, "MongoDB is trying to connect", "error", err.Error(), "attempts left", connAttempts)
		time.Sleep(3 * time.Second)
		connAttempts--
	}

	if m.Cli() == prev {
		return fmt.Errorf("MongoDB: connection failed")
	}
	return nil
}

func (m *Mongo) conn(cliOpts *options.ClientOptions) error {
	newCli, err := mongo.Connect(cliOpts)
	if err != nil {
		return err
	}
	m.setCli(newCli)
	return nil
}

func (m *Mongo) Ping(ctx context.Context) error {
	return m.Cli().Ping(ctx, readpref.PrimaryPreferred())
}

func (m *Mongo) Close(ctx context.Context) {
	m.logger.Info(ctx, "Closing MongoDB")
	if err := m.Cli().Disconnect(ctx); err != nil {
		m.logger.Error(ctx, "Error while closing MongoDB", "error", err.Error())
	}
}
