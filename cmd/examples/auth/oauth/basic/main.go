package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/drummerzzz/goxios"
)

func main() {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token": "test-token", "token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer authServer.Close()

	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
			TokenURL:      authServer.URL,
			ClientID:      "my-client-id",
			ClientSecret:  "my-very-secret-key",
			Scopes:        []string{"read", "write"},
			RefreshBefore: 30 * time.Second,
		}),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	fmt.Println("--- Exemplo de OAuth2 Client Credentials ---")
	fmt.Println("O goxios gerenciará automaticamente a obtenção e renovação do Bearer Token.")

	// Na primeira requisição, o client tentará obter o token no TokenURL configurado
	resp, err := client.Get("/headers").Do(context.Background())
	if err != nil {
		fmt.Printf("Nota: Requisição falhou como esperado (TokenURL fictício): %v\n", err)
		return
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Token obtido e aplicado na requisição.")
	}
}
