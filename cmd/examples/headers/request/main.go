package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	fmt.Println("Fazendo requisição com headers específicos (apenas para esta chamada)...")

	// Adicionando headers diretamente na construção da request
	resp, err := client.Get("/headers").
		Header("X-Request-ID", "req-12345").
		Headers(map[string]string{
			"X-Custom-Info": "specific-data",
			"Cache-Control": "no-cache",
		}).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		data, _ := resp.JsonMap()
		headers := data["headers"].(map[string]any)
		// O httpbin pode normalizar o casing dos headers na resposta
		for k, v := range headers {
			if k == "X-Request-Id" || k == "X-Request-ID" {
				fmt.Printf("  X-Request-ID: %v\n", v)
			}
			if k == "X-Custom-Info" {
				fmt.Printf("  X-Custom-Info: %v\n", v)
			}
		}
	}
}
