package redis

import (
	"context"
	"time"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	SetString(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	GetString(ctx context.Context, key string) ([]byte, error)

	DelKeys(ctx context.Context, keys ...string) error
}

type redisRepo struct {
	logger *log.Logger
	client *redis.Client
}

func NewRedisRespository(cfg *config.Config) *redisRepo {
	opts := redis.Options{
		Addr:            cfg.Redis.Host + ":" + cfg.Redis.Port,
		ClientName:      cfg.Redis.ClientName,
		Username:        cfg.Redis.Username,
		Password:        cfg.Redis.Password,
		MaxRetries:      cfg.Redis.MaxRetries,
		PoolSize:        cfg.Redis.PoolSize,
		MaxIdleConns:    cfg.Redis.MaxIdleConns,
		MaxActiveConns:  cfg.Redis.MaxActiveConns,
		ConnMaxIdleTime: time.Duration(cfg.Redis.MaxIdleConns) * time.Minute,
		ConnMaxLifetime: time.Duration(cfg.Redis.MaxLifeTime) * time.Minute,
	}
	logger := log.With("service", "redis")
	redis.SetLogger(&redisLogger{logger})

	return &redisRepo{
		logger: logger,
		client: redis.NewClient(&opts),
	}
}

func (r *redisRepo) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisRepo) Close(ctx context.Context) {
	r.logger.Info(ctx, "Closing Redis")
	if err := r.client.Close(); err != nil {
		r.logger.Error(ctx, "Error while closing Redis", "error", err.Error())
	}
}

func (r *redisRepo) SetString(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *redisRepo) GetString(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *redisRepo) DelKeys(ctx context.Context, key ...string) error {
	return r.client.Del(ctx, key...).Err()
}
