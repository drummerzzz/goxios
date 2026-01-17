package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

type HttpbinGet struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
}

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	// JsonAs[T] é a forma mais elegante de decodificar em uma struct
	data, err := goxios.JsonAs[HttpbinGet](resp)
	if err != nil {
		log.Fatalf("Erro ao decodificar: %v", err)
	}

	fmt.Println("Dados decodificados com JsonAs:")
	fmt.Printf("  URL: %s\n", data.URL)
	fmt.Printf("  Origin: %s\n", data.Origin)
	fmt.Printf("  User-Agent: %s\n", data.Headers["User-Agent"])
}
