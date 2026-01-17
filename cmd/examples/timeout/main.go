package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Cliente com timeout de 2 segundos
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithTimeout(2*time.Second),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	fmt.Println("--- Teste de Timeout (esperado falha) ---")
	// Forçando um delay de 5 segundos no httpbin para estourar o timeout de 2s
	_, err = client.Get("/delay/5").Do(context.Background())
	if err != nil {
		fmt.Printf("Sucesso: Capturamos o timeout esperado: %v\n", err)
	} else {
		fmt.Println("Erro: O timeout não funcionou como esperado.")
	}

	fmt.Println("\n--- Teste de Sucesso (dentro do timeout) ---")
	resp, err := client.Get("/delay/1").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro inesperado: %v", err)
	}
	if resp.Ok() {
		fmt.Println("Sucesso: Resposta recebida dentro do tempo limite.")
	}
}
