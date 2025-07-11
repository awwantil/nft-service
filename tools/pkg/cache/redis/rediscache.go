package rediscache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"main/tools/pkg/cache"
	coreconfig "main/tools/pkg/core_config"
	"main/tools/pkg/database"
)

// NewRedisClient creates a new Redis client based on the provided configuration.
func NewRedisClient(ctx context.Context, cfg coreconfig.Redis) (cache.CacheClient, error) {

	// подключение к БД
	rdb, err := database.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// возвращаем кеш клиент на redis DB
	return &redisCache{
		client: rdb,
	}, nil
}

// redisCache внутренний клиент для работы с redis DB
type redisCache struct {
	client *redis.Client
}

// Get имплементация метода интерфейса
func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	return data, err
}

// Set имплементация метода интерфейса
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	return err
}

// Exists имплементация метода интерфейса
func (r *redisCache) Exists(ctx context.Context, keys ...string) (uint64, error) {
	exists, err := r.client.Exists(ctx, keys...).Uint64()
	return exists, err
}

// Del имплементация метода интерфейса
func (r *redisCache) Del(ctx context.Context, keys ...string) (uint64, error) {
	res, err := r.client.Del(ctx, keys...).Uint64()
	return res, err
}

// Close имплементация метода интерфейса
func (r *redisCache) Close() error {
	return r.client.Close()
}
