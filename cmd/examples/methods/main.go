package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	ctx := context.Background()

	fmt.Println("--- Exemplos de Métodos HTTP ---")

	// 1. GET
	fmt.Println("\n1. GET:")
	respGet, _ := client.Get("/get").Do(ctx)
	if respGet.Ok() {
		fmt.Println("  [OK] GET realizado com sucesso.")
	}

	// 2. POST
	fmt.Println("\n2. POST:")
	bodyPost := []byte(`{"message": "criando recurso"}`)
	respPost, _ := client.Post("/post", bodyPost).Do(ctx)
	if respPost.Ok() {
		fmt.Println("  [OK] POST realizado com sucesso.")
	}

	// 3. PUT
	fmt.Println("\n3. PUT:")
	bodyPut := []byte(`{"message": "atualizando recurso completo"}`)
	respPut, _ := client.Put("/put", bodyPut).Do(ctx)
	if respPut.Ok() {
		fmt.Println("  [OK] PUT realizado com sucesso.")
	}

	// 4. PATCH
	fmt.Println("\n4. PATCH:")
	bodyPatch := []byte(`{"message": "atualização parcial"}`)
	respPatch, _ := client.Patch("/patch", bodyPatch).Do(ctx)
	if respPatch.Ok() {
		fmt.Println("  [OK] PATCH realizado com sucesso.")
	}

	// 5. DELETE
	fmt.Println("\n5. DELETE:")
	respDelete, _ := client.Delete("/delete").Do(ctx)
	if respDelete.Ok() {
		fmt.Println("  [OK] DELETE realizado com sucesso.")
	}

	// 6. HEAD
	fmt.Println("\n6. HEAD:")
	respHead, _ := client.Head("/get").Do(ctx)
	if respHead.Ok() {
		fmt.Println("  [OK] HEAD realizado com sucesso.")
		fmt.Printf("  Content-Length: %s\n", respHead.Header.Get("Content-Length"))
	}

	// 7. OPTIONS
	fmt.Println("\n7. OPTIONS:")
	respOptions, _ := client.Options("/get").Do(ctx)
	if respOptions.Ok() {
		fmt.Println("  [OK] OPTIONS realizado com sucesso.")
		fmt.Printf("  Allow: %s\n", respOptions.Header.Get("Allow"))
	}
}

