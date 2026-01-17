package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implementa a interface TokenCache usando Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache cria uma nova instância de RedisCache.
func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

