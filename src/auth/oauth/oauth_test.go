package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenSource_RenewToken(t *testing.T) {
	t.Parallel()

	var nowUnix atomic.Int64
	nowUnix.Store(time.Unix(1000, 0).Unix())
	now := func() time.Time { return time.Unix(nowUnix.Load(), 0) }

	var tokenCalls atomic.Int64
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		tok := fmt.Sprintf("t%d", tokenCalls.Load())
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": tok,
			"token_type":   "Bearer",
			"expires_in":   1,
		})
	}))
	t.Cleanup(tokenSrv.Close)

	src := NewTokenSource(http.DefaultClient, nil, Config[DefaultTokenResponse]{
		TokenURL:      tokenSrv.URL,
		ClientID:      "id",
		ClientSecret:  "secret",
		RefreshBefore: 0,
		Now:           now,
	})

	ctx := context.Background()

	// 1st call: fetch token
	tok1, err := src.Token(ctx)
	if err != nil {
		t.Fatalf("Token() err=%v", err)
	}
	if tok1 != "t1" {
		t.Fatalf("expected t1; got=%v", tok1)
	}

	// Advance time to expire and force refresh
	nowUnix.Add(2)

	tok2, err := src.Token(ctx)
	if err != nil {
		t.Fatalf("Token() err=%v", err)
	}
	if tok2 != "t2" {
		t.Fatalf("expected t2; got=%v", tok2)
	}
	if tokenCalls.Load() != 2 {
		t.Fatalf("expected 2 calls to token endpoint; got=%d", tokenCalls.Load())
	}
}

type memTokenCache struct {
	m map[string]string
}

func (c *memTokenCache) Get(ctx context.Context, key string) (string, error) {
	if c == nil || c.m == nil {
		return "", nil
	}
	return c.m[key], nil
}

func (c *memTokenCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c.m == nil {
		c.m = map[string]string{}
	}
	c.m[key] = value
	return nil
}

func TestTokenSource_UsesExternalCache(t *testing.T) {
	t.Parallel()

	cache := &memTokenCache{}

	var nowUnix atomic.Int64
	nowUnix.Store(time.Now().Unix())
	now := func() time.Time { return time.Unix(nowUnix.Load(), 0) }

	var tokenCalls atomic.Int64
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "cached-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	t.Cleanup(tokenSrv.Close)

	src := NewTokenSource(http.DefaultClient, nil, Config[DefaultTokenResponse]{
		TokenURL:      tokenSrv.URL,
		ClientID:      "id",
		ClientSecret:  "secret",
		Cache:         cache,
		RefreshBefore: 30 * time.Second,
		Now:           now,
	})

	// 1st request: will call token endpoint and write to cache
	tok1, err := src.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() err=%v", err)
	}
	if tok1 != "cached-token" {
		t.Fatalf("unexpected token: %v", tok1)
	}
	if tokenCalls.Load() != 1 {
		t.Fatalf("expected 1 call to token endpoint; got=%d", tokenCalls.Load())
	}

	// New Source, same config/cache: should read from cache and not call token endpoint
	src2 := NewTokenSource(http.DefaultClient, nil, Config[DefaultTokenResponse]{
		TokenURL:      tokenSrv.URL,
		ClientID:      "id",
		ClientSecret:  "secret",
		Cache:         cache,
		RefreshBefore: 30 * time.Second,
		Now:           now,
	})

	tok3, err := src2.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() err=%v", err)
	}
	if tok3 != "cached-token" {
		t.Fatalf("unexpected token: %v", tok3)
	}
	if tokenCalls.Load() != 1 {
		t.Fatalf("new Source should use cache and not call token endpoint; got=%d", tokenCalls.Load())
	}
}

type customTokenResponse struct {
	MyToken   string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func (r customTokenResponse) GetAccessToken() string { return r.MyToken }
func (r customTokenResponse) GetExpiresIn() int64   { return r.ExpiresAt }

func TestTokenSource_CustomResponse(t *testing.T) {
	t.Parallel()

	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"token":      "custom-token",
			"expires_at": 123,
		})
	}))
	t.Cleanup(tokenSrv.Close)

	src := NewTokenSource(http.DefaultClient, nil, Config[customTokenResponse]{
		TokenURL:     tokenSrv.URL,
		ClientID:     "id",
		ClientSecret: "secret",
	})

	tok, err := src.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() err=%v", err)
	}
	if tok != "custom-token" {
		t.Fatalf("expected custom-token; got=%v", tok)
	}
}

