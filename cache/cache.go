package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) *Result
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Do(ctx context.Context, args ...any) (*Result, error)
}
type Result struct {
	val   string
	error error
}

// 获取int
func (r *Result) Int() (int, error) {
	return 0, nil
}

// 获取结果
func (r *Result) Val() string {
	return r.val
}

// 获取错误
func (r *Result) Err() error {
	return r.error
}
