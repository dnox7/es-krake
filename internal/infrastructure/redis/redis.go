package redis

import (
	"context"
	"sync"
	"time"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	mu     sync.RWMutex
	logger *log.Logger
	client *redis.Client
	params redisParams
}

type redisParams struct {
	address         string
	clientName      string
	maxRetries      int
	poolSize        int
	maxIdleConns    int
	maxActiveConns  int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

var (
	redisInstance *RedisRepo
	once          sync.Once
)

func NewRedisRespository(ctx context.Context, cfg *config.Config, cred *config.RedisCredentials) *RedisRepo {
	once.Do(func() {
		repo, err := initRedis(ctx, cfg, cred)
		if err != nil {
			panic(err)
		}
		redisInstance = repo
	})
	return redisInstance
}

func initRedis(ctx context.Context, cfg *config.Config, cred *config.RedisCredentials) (*RedisRepo, error) {
	repo := &RedisRepo{
		logger: log.With("service", "redis"),
		params: redisParams{
			address:         cfg.Redis.Host + ":" + cfg.Redis.Port,
			clientName:      cfg.Redis.ClientName,
			maxRetries:      cfg.Redis.MaxRetries,
			poolSize:        cfg.Redis.PoolSize,
			maxIdleConns:    cfg.Redis.MaxIdleConns,
			maxActiveConns:  cfg.Redis.MaxActiveConns,
			connMaxIdleTime: time.Duration(cfg.Redis.MaxIdleConns) * time.Minute,
			connMaxLifetime: time.Duration(cfg.Redis.MaxLifeTime) * time.Minute,
		},
	}
	redis.SetLogger(&redisLogger{repo.logger})
	if err := repo.RetryConn(ctx, cred); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RedisRepo) RetryConn(ctx context.Context, cred *config.RedisCredentials) error {
	opts := redis.Options{
		Addr:            r.params.address,
		ClientName:      r.params.clientName,
		Username:        cred.Username,
		Password:        cred.Password,
		MaxRetries:      r.params.maxRetries,
		PoolSize:        r.params.poolSize,
		MaxActiveConns:  r.params.maxActiveConns,
		MaxIdleConns:    r.params.maxIdleConns,
		ConnMaxLifetime: r.params.connMaxLifetime,
		ConnMaxIdleTime: r.params.connMaxIdleTime,
	}
	client := redis.NewClient(&opts)
	r.setCli(client)

	if err := r.Ping(ctx); err != nil {
		r.logger.Error(ctx, "error while pinging Redis", "error", err)
		return err
	}
	return nil
}

func (r *RedisRepo) cli() *redis.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

func (r *RedisRepo) setCli(newCli *redis.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = newCli
}

func (r *RedisRepo) Ping(ctx context.Context) error {
	return r.cli().Ping(ctx).Err()
}

func (r *RedisRepo) Close(ctx context.Context) {
	r.logger.Info(ctx, "Closing Redis")
	if err := r.cli().Close(); err != nil {
		r.logger.Error(ctx, "Error while closing Redis", "error", err.Error())
	}
}
