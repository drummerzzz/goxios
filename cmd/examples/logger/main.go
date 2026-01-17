package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drummerzzz/goxios"
	"go.uber.org/zap"
)

func main() {
	// 1. Criando um logger do Zap
	// O goxios utiliza o zap.Logger para logging interno.
	// NewDevelopment é ótimo para ver logs detalhados no console.
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // garante que todos os logs sejam escritos antes de sair

	// 2. Configurando o cliente com o logger
	client, err := goxios.New(
		goxios.WithBaseURL("https://httpbin.org"),
		goxios.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Println("--- Exemplo de Logging com Zap ---")
	fmt.Println("O goxios logará detalhes da requisição e resposta no console.")

	// 3. Fazendo uma requisição
	// O logger irá capturar o início, erros (se houver) e o JSON da resposta.
	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		log.Fatalf("Erro na requisição: %v", err)
	}

	if resp.Ok() {
		// Ao chamar métodos de leitura de JSON, o goxios loga o conteúdo do body se o logger estiver ativo.
		_, _ = resp.Json()
		fmt.Println("\nRequisição completada. Verifique os logs acima para detalhes do body e status.")
	}
}
