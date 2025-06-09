package redis

import (
	"context"
	"time"
)

type RedisRepository interface {
	SetString(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	GetString(ctx context.Context, key string) ([]byte, error)

	PushBackList(ctx context.Context, key string, ttl time.Duration, vals ...interface{}) error
	GetRangeList(ctx context.Context, key string, start, end int64, container interface{}) error
	PopBackList(ctx context.Context, key string, count int) error

	AddSet(ctx context.Context, key string, vals ...interface{}) error
	GetSet(ctx context.Context, key string, container interface{}) error
	RemoveEleSet(ctx context.Context, key string, vals ...interface{}) error

	DelKeys(ctx context.Context, keys ...string) error
}

func (r *RedisRepo) SetString(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return r.cli().Set(ctx, key, val, ttl).Err()
}

func (r *RedisRepo) GetString(ctx context.Context, key string) ([]byte, error) {
	return r.cli().Get(ctx, key).Bytes()
}

func (r *RedisRepo) PushBackList(ctx context.Context, key string, ttl time.Duration, vals ...interface{}) error {
	err := r.cli().RPush(ctx, key, vals...).Err()
	if err != nil {
		return err
	}
	return r.cli().Expire(ctx, key, ttl).Err()
}

func (r *RedisRepo) GetRangeList(ctx context.Context, key string, start, end int64, container interface{}) error {
	res := r.cli().LRange(ctx, key, start, end)
	if res.Err() != nil {
		return res.Err()
	}
	return res.ScanSlice(container)
}

func (r *RedisRepo) PopBackList(ctx context.Context, key string, count int) error {
	return r.cli().RPopCount(ctx, key, count).Err()
}

func (r *RedisRepo) AddSet(ctx context.Context, key string, vals ...interface{}) error {
	return r.cli().SAdd(ctx, key, vals...).Err()
}

func (r *RedisRepo) GetSet(ctx context.Context, key string, container interface{}) error {
	res := r.cli().SMembers(ctx, key)
	if res.Err() != nil {
		return res.Err()
	}
	return res.ScanSlice(container)
}

func (r *RedisRepo) RemoveEleSet(ctx context.Context, key string, vals ...interface{}) error {
	return r.cli().SRem(ctx, key, vals...).Err()
}

func (r *RedisRepo) DelKeys(ctx context.Context, keys ...string) error {
	return r.cli().Del(ctx, keys...).Err()
}
