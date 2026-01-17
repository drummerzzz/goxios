package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/drummerzzz/goxios"
)

// MemoryCache é uma implementação simples da interface TokenCache
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.store[key], nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
	return nil
}

func main() {
	// Servidor de autenticação fake
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token": "token-cacheado-123", "token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer authServer.Close()

	// 1. Instanciando nosso cache customizado
	myCache := &MemoryCache{store: make(map[string]string)}

	// 2. Configurando o cliente com OAuth e o Cache Externo
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
			TokenURL:     authServer.URL,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Cache:        myCache, // <-- Plugando o cache aqui
		}),
	)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("--- Exemplo de Cache Customizado para OAuth ---")

	// Primeira requisição: vai buscar no servidor de auth e salvar no cache
	fmt.Println("1. Fazendo primeira requisição (deve popular o cache)...")
	_, _ = client.Get("/headers").Do(context.Background())

	// Verificando se algo foi salvo no cache (qualquer chave)
	myCache.mu.RLock()
	fmt.Printf("Itens no cache: %d\n", len(myCache.store))
	myCache.mu.RUnlock()

	// Segunda requisição: vai ler do cache diretamente
	fmt.Println("\n2. Fazendo segunda requisição (deve ler do cache)...")
	resp, _ := client.Get("/headers").Do(context.Background())

	if resp.Ok() {
		fmt.Println("Sucesso! Token recuperado e utilizado via cache.")
	}
}
