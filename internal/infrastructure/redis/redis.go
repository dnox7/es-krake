package redis

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"

	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
}

type redisService struct {
	logger *log.Logger
	client *redis.Client
}

func NewRedisService(cfg *config.Config) CacheStore {
	opts := redis.Options{
		// Addr:            cfg.Redis.Host + ":" + cfg.Redis.Port,
		// ClientName:      cfg.Redis.ClientName,
		// Username:        cfg.Redis.Username,
		// Password:        cfg.Redis.Password,
		// MaxRetries:      cfg.Redis.MaxRetries,
		// PoolSize:        cfg.Redis.PoolSize,
		// MaxIdleConns:    cfg.Redis.MaxIdleConns,
		// MaxActiveConns:  cfg.Redis.MaxActiveConns,
		// ConnMaxIdleTime: time.Duration(cfg.Redis.MaxIdleConns) * time.Minute,
		// ConnMaxLifetime: time.Duration(cfg.Redis.MaxLifeTime) * time.Minute,
	}
	return &redisService{
		logger: log.With("service", "redis"),
		client: redis.NewClient(&opts),
	}
}
