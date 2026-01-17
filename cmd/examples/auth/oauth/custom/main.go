package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/drummerzzz/goxios"
	"github.com/drummerzzz/goxios/src/auth/oauth"
)

// MyCustomTokenResponse define uma estrutura de resposta OAuth não convencional.
// Alguns clientes podem retornar campos como "token" em vez de "access_token"
// ou "expires_at" em vez de "expires_in".
type MyCustomTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// GetAccessToken implementa a interface oauth.TokenResponse.
func (r MyCustomTokenResponse) GetAccessToken() string {
	return r.Token
}

// GetExpiresIn implementa a interface oauth.TokenResponse.
func (r MyCustomTokenResponse) GetExpiresIn() int64 {
	return r.ExpiresAt
}

func main() {
	// 1. Simula um servidor de autenticação que retorna um formato customizado.
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Retorna "token" e "expires_at" em vez do padrão OAuth2.
		w.Write([]byte(`{
			"token": "custom-provided-token-123",
			"expires_at": 3600
		}`))
	}))
	defer authServer.Close()

	// 2. Configura o goxios usando WithOAuthClientCredentialsCustom com a struct genérica.
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithOAuthClientCredentialsCustom(oauth.Config[MyCustomTokenResponse]{
			TokenURL:      authServer.URL,
			ClientID:      "custom-client",
			ClientSecret:  "custom-secret",
			RefreshBefore: 30 * time.Second,
		}),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	fmt.Println("--- Exemplo de OAuth2 com Resposta Customizada (Generics) ---")
	fmt.Println("Usando uma struct que mapeia 'token' e 'expires_at' em vez do padrão.")

	// Na primeira requisição, o client obterá o token usando a struct MyCustomTokenResponse.
	resp, err := client.Get("/headers").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro na requisição: %v", err)
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Token customizado obtido e aplicado corretamente!")
	} else {
		fmt.Printf("Falha na requisição: %d\n", resp.StatusCode)
	}
}
