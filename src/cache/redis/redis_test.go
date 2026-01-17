package redis

import (
	"context"
	"testing"
)

func TestNewRedisCache(t *testing.T) {
	addr := "localhost:6379"
	cache := NewRedisCache(addr)

	if cache == nil {
		t.Fatal("expected non-nil RedisCache")
	}

	if cache.client == nil {
		t.Fatal("expected non-nil redis.Client")
	}
}

func TestRedisCache_Offline(t *testing.T) {
	// Test behavior when Redis is offline (or invalid address)
	cache := NewRedisCache("localhost:1") // Unlikely port to be open
	ctx := context.Background()

	_, err := cache.Get(ctx, "test-key")
	if err == nil {
		t.Error("expected error when trying Get on offline redis")
	}

	err = cache.Set(ctx, "test-key", "value", 0)
	if err == nil {
		t.Error("expected error when trying Set on offline redis")
	}
}

