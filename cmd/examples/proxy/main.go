package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/drummerzzz/goxios"
)

func main() {
	// Exemplo de como configurar um proxy HTTP
	// Criamos um servidor de proxy local para demonstração
	proxyServer := createProxyServer()
	defer proxyServer.Close()

	client, err := goxios.New(
		goxios.WithBaseURL("http://httpbin.org"),
		goxios.WithProxyURL(proxyServer.URL),
	)
	if err != nil {
		log.Fatalf("Erro ao criar cliente: %v", err)
	}

	fmt.Println("--- Exemplo de Configuração de Proxy ---")
	fmt.Printf("Tentando acessar via proxy: %s\n", proxyServer.URL)
	fmt.Println("(Este exemplo pode falhar se você não tiver um proxy ativo)")

	resp, err := client.Get("/get").Do(context.Background())
	if err != nil {
		fmt.Printf("Nota: Requisição falhou (provavelmente proxy offline): %v\n", err)
		return
	}

	if resp.Ok() {
		fmt.Println("Sucesso: Requisição via proxy completada.")
	}
}

func createProxyServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)

		// copia body
		io.Copy(w, resp.Body)
	}))
}
