package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
	"go.uber.org/zap"
)

func main() {
	user := "admin"
	pass := "secret123"

	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	client, err := goxios.New(
		goxios.WithLogger(logger),
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithBasicAuth(user, pass),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	fmt.Printf("--- Teste de BasicAuth para usuário '%s' ---\n", user)

	// Endpoint do httpbin que valida BasicAuth
	url := fmt.Sprintf("/basic-auth/%s/%s", user, pass)
	resp, err := client.Get(url).Do(context.Background())
	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Autenticação básica funcionou!")
		goxios.JsonAs[map[string]any](resp)
	} else {
		fmt.Printf("Erro: Falha na autenticação. Status: %d\n", resp.StatusCode)
	}
}
