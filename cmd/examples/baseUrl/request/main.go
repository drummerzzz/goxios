package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Cliente configurado com um BaseURL padrão
	client, _ := goxios.New()
	resp, err := client.Get("https://httpbin.org/get").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		data, _ := resp.JsonMap()
		fmt.Printf("  URL chamada: %v\n", data["url"])
	}
}
