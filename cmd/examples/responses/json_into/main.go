package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

type HttpbinGet struct {
	URL string `json:"url"`
}

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, _ := client.Get("/get").Do(context.Background())

	// JsonInto permite decodificar em uma struct já existente
	var data HttpbinGet
	err := resp.JsonInto(&data)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("Dados decodificados com JsonInto:")
	fmt.Printf("  URL: %s\n", data.URL)
}
