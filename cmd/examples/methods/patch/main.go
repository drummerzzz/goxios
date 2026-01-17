package main

import (
	"context"
	"fmt"
	"log"

	goxios "github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	body := []byte(`{"status": "partially-updated"}`)
	resp, err := client.Patch("/patch", body).Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		fmt.Println("PATCH executado com sucesso!")
		data, _ := resp.JsonMap()
		fmt.Printf("Campos alterados: %v\n", data["json"])
	}
}

