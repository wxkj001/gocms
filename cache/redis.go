package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	con *redis.Client
}

func NewRedisCache(con *redis.Client) *RedisCache {
	return &RedisCache{con: con}
}
func (r *RedisCache) Get(ctx context.Context, key string) *Result {
	res, err := r.con.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return &Result{val: "", error: nil}
		}
		return &Result{val: "", error: err}
	}
	return &Result{val: res, error: nil}
}
func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	return r.con.Set(ctx, key, val, expiration).Err()
}
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.con.Del(ctx, keys...).Err()
}
func (r *RedisCache) Do(ctx context.Context, args ...any) (*Result, error) {
	res, err := r.con.Do(ctx, args).Text()
	if err != nil {
		return nil, err
	}
	return &Result{val: res}, nil
}
