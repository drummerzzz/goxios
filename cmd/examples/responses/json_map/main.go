package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	// JsonMap decodifica o JSON em um map[string]any
	data, err := resp.JsonMap()
	if err != nil {
		log.Fatalf("Erro ao ler mapa: %v", err)
	}

	fmt.Println("Dados lidos via JsonMap:")
	fmt.Printf("  URL: %v\n", data["url"])
	if headers, ok := data["headers"].(map[string]any); ok {
		fmt.Printf("  User-Agent: %v\n", headers["User-Agent"])
	}
}
