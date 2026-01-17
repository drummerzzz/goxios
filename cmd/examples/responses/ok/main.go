package main

import (
	"context"
	"fmt"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, _ := goxios.New(goxios.WithBaseURL("https://httpbin.org"))

	// 200 OK
	resp200, _ := client.Get("/status/200").Do(context.Background())
	fmt.Printf("Status 200 - Ok(): %v\n", resp200.Ok())

	// 404 Not Found
	resp404, _ := client.Get("/status/404").Do(context.Background())
	fmt.Printf("Status 404 - Ok(): %v\n", resp404.Ok())

	// 500 Internal Server Error
	resp500, _ := client.Get("/status/500").Do(context.Background())
	fmt.Printf("Status 500 - Ok(): %v\n", resp500.Ok())
}
