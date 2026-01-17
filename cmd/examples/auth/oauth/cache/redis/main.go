package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
	"github.com/drummerzzz/goxios/src/cache/redis"
)

func main() {
	// 1. Criando o cache do Redis (interno da lib)
	// Assume que existe um Redis rodando em localhost:6379
	redisCache := redis.NewRedisCache("localhost:6379")

	// 2. Configurando o cliente com OAuth e Redis
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
			TokenURL:     "https://sua-api.com/oauth/token",
			ClientID:     "meu-id",
			ClientSecret: "meu-secret",
			Cache:        redisCache, // <-- Usando Redis para compartilhar tokens entre instâncias
		}),
	)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("--- Exemplo de Cache com Redis para OAuth ---")
	fmt.Println("Configuração aplicada. Ao fazer requisições, o token será persistido no Redis.")
	fmt.Println("(Este exemplo requer um Redis ativo para funcionar na prática)")

	// O uso do client permanece transparente
	resp, err := client.Get("/headers").Do(context.Background())
	if err != nil {
		fmt.Printf("Nota: Requisição falhou como esperado (Redis ou URL de auth offline): %v\n", err)
		return
	}

	if resp.Ok() {
		fmt.Println("Requisição concluída com sucesso usando token do Redis.")
	}
}
