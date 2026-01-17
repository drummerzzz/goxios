package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	// O BaseURL é configurado uma vez no cliente.
	// Todas as requisições usarão este prefixo se um caminho relativo for passado.
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
	)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("Usando BaseURL padrão do cliente...")

	// Passando apenas o caminho relativo "/get"
	resp, _ := client.Get("/get").Do(context.Background())

	if resp.Ok() {
		data, _ := resp.JsonMap()
		fmt.Printf("  URL final chamada: %v\n", data["url"])
	}
}
