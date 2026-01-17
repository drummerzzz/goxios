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
go get github.com/drummerzzz/goxios
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
        fmt.Println("Status:", resp.StatusCode())
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

