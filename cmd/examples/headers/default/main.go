package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Configurando headers que estarão presentes em TODAS as requisições deste cliente
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithDefaultHeader("X-App-Name", "GoxiosExample"),
		goxios.WithDefaultHeaders(map[string]string{
			"X-Environment": "Development",
			"Accept":        "application/json",
		}),
	)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("Fazendo requisição com headers padrão do cliente...")
	resp, _ := client.Get("/headers").Do(context.Background())

	if resp.Ok() {
		data, _ := resp.JsonMap()
		headers := data["headers"].(map[string]any)
		fmt.Printf("  X-App-Name: %v\n", headers["X-App-Name"])
		fmt.Printf("  X-Environment: %v\n", headers["X-Environment"])
	}
}
