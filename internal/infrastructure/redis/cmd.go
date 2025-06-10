package redis

import (
	"context"
	"encoding/json"
	"time"
)

type RedisRepository interface {
	SetJSON(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	GetJSON(ctx context.Context, key string, container interface{}) error

	CheckExists(ctx context.Context, key string) (bool, error)
	DelKeys(ctx context.Context, keys ...string) error
}

func (r *RedisRepo) SetJSON(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return r.cli().Set(ctx, key, data, ttl).Err()
}

func (r *RedisRepo) GetJSON(ctx context.Context, key string, container interface{}) error {
	bytes, err := r.cli().Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, container)
}

func (r *RedisRepo) CheckExists(ctx context.Context, key string) (bool, error) {
	res := r.cli().Exists(ctx, key)
	if res.Err() != nil {
		return false, res.Err()
	}
	return res.Val() == 1, nil
}

func (r *RedisRepo) DelKeys(ctx context.Context, keys ...string) error {
	return r.cli().Del(ctx, keys...).Err()
}
