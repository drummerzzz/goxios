package cache

import (
	"context"
	"time"
)

// TokenCache é uma interface mínima para plugar sistemas de cache (ex: Redis, Memcached) para tokens.
type TokenCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}
