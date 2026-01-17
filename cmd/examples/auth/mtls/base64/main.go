package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Carrega os certificados do arquivo
	certBase64, keyBase64 := getCertAndKeyBase64()
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithMTLSFromBase64(certBase64, keyBase64),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente com mTLS (base64): %v", err)
	}

	fmt.Println("--- Exemplo de mTLS usando Strings Base64 ---")

	// Nota: Este exemplo falhará na execução pois as strings acima são inválidas,
	// mas demonstra a sintaxe de configuração.
	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		fmt.Printf("Nota: Requisição falhou como esperado (strings inválidas): %v\n", err)
		return
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Transporte TLS configurado com base64.")
	}
}

func getCertAndKeyBase64() (string, string) {
	certFile := "cmd/examples/auth/mtls/cert.pem"
	keyFile := "cmd/examples/auth/mtls/key.pem"

	certBytes, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Erro ao ler arquivo de chave: %v", err)
	}
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Erro ao criar cliente com mTLS (base64): %v", err)
	}
	certBase64 := base64.StdEncoding.EncodeToString(certBytes)
	keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)
	return certBase64, keyBase64
}
