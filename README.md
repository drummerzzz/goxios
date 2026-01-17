# Goxios

Goxios é uma biblioteca de cliente HTTP para Go (Golang), inspirada no Axios do ecossistema JavaScript. Ela oferece uma interface fluida, organizada e poderosa para realizar requisições HTTP, lidar com autenticação, cache e mTLS de forma simplificada.

## 🚀 Funcionalidades

- **API Fluida**: Construção de requisições de forma encadeada.
- **Cliente Reutilizável**: Configure um cliente uma vez e reuse em toda a aplicação.
- **Autenticação Integrada**: Suporte nativo a Basic Auth, Bearer Token e OAuth2 (Client Credentials).
- **Gerenciamento de Cache**: Cache de tokens OAuth2 extensível (com implementação Redis inclusa).
- **Suporte a mTLS**: Configuração simplificada de certificados via arquivos ou strings Base64.
- **Tratamento de JSON**: Facilidade para decodificar respostas JSON diretamente em structs ou maps.
- **Logging**: Integração com a biblioteca `zap` para logs detalhados.
- **Extensível**: Sistema de opções para configurar o cliente de acordo com sua necessidade.

## 📦 Instalação

```bash
go get github.com/drummerzzz/goxios@v0.1.0
```

## 🛠️ Como usar

Para uma documentação detalhada de todas as funcionalidades, consulte a [Documentação da API](API.md).

### Uso Básico

```go
package main

import (
    "context"
    "fmt"
    "github.com/drummerzzz/goxios"
)

func main() {
    client, _ := goxios.New(
        goxios.WithBaseURL("https://api.github.com"),
    )

    resp, err := client.Get("/repos/drummerzzz/goxios").Do(context.Background())
    if err != nil {
        panic(err)
    }

    if resp.Ok() {
        fmt.Println("Status:", resp.StatusCode)
    }
}
```

### Autenticação OAuth2 com Redis

O Goxios facilita o uso de OAuth2, permitindo inclusive o cache de tokens no Redis para compartilhamento entre múltiplas instâncias da sua aplicação.

```go
import (
    "github.com/drummerzzz/goxios"
    "github.com/drummerzzz/goxios/src/cache/redis"
)

func main() {
    redisCache := redis.NewRedisCache("localhost:6379")

    client, _ := goxios.New(
        goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
            TokenURL:     "https://sua-api.com/oauth/token",
            ClientID:     "meu-id",
            ClientSecret: "meu-secret",
            Cache:        redisCache,
        }),
    )
}
```

### Configuração de mTLS

Você pode configurar mTLS globalmente no cliente ou apenas para uma requisição específica.

```go
// No cliente (global)
client, _ := goxios.New(
    goxios.WithMTLSFromFile("cert.pem", "key.pem"),
)

// Ou em uma requisição específica
resp, _ := client.Get("/seguro").
    MTLS(&goxios.Certificate{
        MtlsCertBase64: "...",
        MtlsKeyBase64:  "...",
    }).
    Do(ctx)
```

### Manipulação de JSON

```go
type User struct {
    Name string `json:"name"`
}

// Decodificando em uma struct
user, err := goxios.JsonAs[User](resp)

// Ou decodificando em um map
data, err := resp.JsonMap()
```

## 📖 Exemplos

Você pode encontrar exemplos detalhados de cada funcionalidade no diretório [`cmd/examples`](cmd/examples).

### Autenticação
- [Basic Auth](cmd/examples/auth/basicauth/main.go)
- [OAuth2 Básico](cmd/examples/auth/oauth/basic/main.go)
- [OAuth2 com Resposta Customizada (Generics)](cmd/examples/auth/oauth/custom/main.go)
- [OAuth2 com Cache Customizado](cmd/examples/auth/oauth/cache/custom/main.go)
- [OAuth2 com Cache Redis](cmd/examples/auth/oauth/cache/redis/main.go)
- [mTLS via Base64](cmd/examples/auth/mtls/base64/main.go)
- [mTLS via Arquivo Único](cmd/examples/auth/mtls/file/main.go)
- [mTLS via Múltiplos Arquivos](cmd/examples/auth/mtls/files/main.go)

### Configurações do Cliente
- [Base URL Padrão](cmd/examples/baseUrl/default/main.go)
- [Base URL por Requisição](cmd/examples/baseUrl/request/main.go)
- [Headers Padrão](cmd/examples/headers/default/main.go)
- [Headers por Requisição](cmd/examples/headers/request/main.go)
- [Logger (Zap)](cmd/examples/logger/main.go)
- [Proxy](cmd/examples/proxy/main.go)
- [Timeout](cmd/examples/timeout/main.go)

### Métodos HTTP
- [GET](cmd/examples/methods/get/main.go)
- [POST](cmd/examples/methods/post/main.go)
- [PUT](cmd/examples/methods/put/main.go)
- [DELETE](cmd/examples/methods/delete/main.go)
- [PATCH](cmd/examples/methods/patch/main.go)
- [HEAD](cmd/examples/methods/head/main.go)
- [OPTIONS](cmd/examples/methods/options/main.go)

### Respostas e JSON
- [Decodificar como Struct (Generics)](cmd/examples/responses/json_as/main.go)
- [Decodificar para Struct Existente](cmd/examples/responses/json_into/main.go)
- [Decodificar como Map](cmd/examples/responses/json_map/main.go)
- [Obter Headers da Resposta](cmd/examples/responses/headers/main.go)
- [Obter Bytes Brutos](cmd/examples/responses/raw_bytes/main.go)
- [Verificar Sucesso (OK)](cmd/examples/responses/ok/main.go)

## 📂 Estrutura do Projeto

O projeto segue uma organização moderna e limpa:

- `src/`: Contém o código principal da biblioteca.
  - `auth/`: Implementações de autenticação (OAuth2, etc).
  - `cache/`: Interfaces e implementadores de cache (Redis).
  - `request/`: Lógica de construção e execução de requisições.
  - `response/`: Wrapper para respostas HTTP e helpers JSON.
- `internal/`: Pacotes utilitários de uso interno (TLS, etc), acessíveis por todo o módulo.
- `cmd/examples/`: Exemplos práticos de todas as funcionalidades.

## 🧪 Testes

Para rodar os testes:

```bash
go test ./...
```

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

