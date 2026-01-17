package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Caminhos para os arquivos de certificado e chave
	// Nota: Certifique-se que os arquivos existam no caminho especificado
	certFile := "cmd/examples/auth/mtls/cert.pem"
	keyFile := "cmd/examples/auth/mtls/key.pem"

	// Criando o cliente configurado para mTLS usando arquivos
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithMTLSFromFile(certFile, keyFile),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente com mTLS (arquivo): %v", err)
	}

	fmt.Println("--- Exemplo de mTLS usando Arquivos ---")

	// Nota: httpbin.org não valida o certificado, mas isso demonstra
	// que o transporte TLS foi configurado corretamente sem erros de leitura.
	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		fmt.Printf("Nota: Requisição falhou (provavelmente certificados ausentes): %v\n", err)
		return
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Transporte TLS configurado com arquivos de certificado.")
	}
}

