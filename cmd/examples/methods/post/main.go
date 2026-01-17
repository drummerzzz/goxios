package main

import (
	"context"
	"fmt"
	"log"

	goxios "github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	body := []byte(`{"name": "Goxios", "type": "Library"}`)
	resp, err := client.Post("/post", body).Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		fmt.Println("POST executado com sucesso!")
		data, _ := resp.JsonMap()
		fmt.Printf("Dados enviados: %v\n", data["json"])
	}
}


