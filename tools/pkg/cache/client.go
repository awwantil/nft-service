package cache

import (
	"context"
	"time"
)

// CacheClient интерфейс для работы с кешем
type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Exists(ctx context.Context, keys ...string) (uint64, error)
	Del(ctx context.Context, keys ...string) (uint64, error)
	Close() error
}
