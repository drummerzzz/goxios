package main

import (
	"context"
	"fmt"
	"log"

	goxios "github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	body := []byte(`{"id": 1, "status": "updated"}`)
	resp, err := client.Put("/put", body).Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		fmt.Println("PUT executado com sucesso!")
		data, _ := resp.JsonMap()
		fmt.Printf("Dados atualizados: %v\n", data["json"])
	}
}

