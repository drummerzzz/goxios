package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, _ := client.Get("/get").Do(context.Background())

	// Json() (apesar do nome) retorna os bytes crus do body e faz cache
	bodyBytes, err := resp.Json()
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("Body recebido (Raw Bytes):")
	fmt.Printf("  Tamanho: %d bytes\n", len(bodyBytes))
	fmt.Printf("  Conteúdo (primeiros 50 caracteres): %s...\n", string(bodyBytes[:50]))
}
