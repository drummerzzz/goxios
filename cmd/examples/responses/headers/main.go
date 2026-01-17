package main

import (
	"context"
	"fmt"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	resp, _ := client.Get("/get").Do(context.Background())

	fmt.Println("Lendo headers da resposta:")
	fmt.Printf("  Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("  Server: %s\n", resp.Header.Get("Server"))
	fmt.Printf("  Date: %s\n", resp.Header.Get("Date"))
}
