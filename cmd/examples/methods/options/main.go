package main

import (
	"context"
	"fmt"
	"log"

	goxios "github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, err := client.Options("/get").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	if resp.Ok() {
		fmt.Println("OPTIONS executado com sucesso!")
		fmt.Printf("Métodos permitidos: %s\n", resp.Header.Get("Allow"))
	}
}


