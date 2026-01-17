# Documentação da API Goxios

Esta documentação detalha as funcionalidades da biblioteca Goxios, utilizando os exemplos contidos em `cmd/examples` como base.

## 1. Criação e Configuração do Cliente

O cliente é criado usando a função `goxios.New`, que aceita diversas opções de configuração.

### Base URL
Define um prefixo para todas as requisições feitas por esse cliente. ([Exemplos](cmd/examples/baseUrl))
```go
client, _ := goxios.New(goxios.WithBaseURL("https://api.example.com"))
// Requisições agora podem usar caminhos relativos:
client.Get("/users").Do(ctx)
```

### Timeout
Define o tempo limite para as requisições. ([Exemplo](cmd/examples/timeout))
```go
client, _ := goxios.New(goxios.WithTimeout(10 * time.Second))
```

### Headers Padrão
Define headers que serão enviados em todas as requisições. ([Exemplos](cmd/examples/headers))
```go
client, _ := goxios.New(
    goxios.WithDefaultHeader("X-App-Name", "MinhaApp"),
    goxios.WithDefaultHeaders(map[string]string{
        "Accept": "application/json",
    }),
)
```

### Logging
Integração com o `zap.Logger`. ([Exemplo](cmd/examples/logger))
```go
logger, _ := zap.NewDevelopment()
client, _ := goxios.New(goxios.WithLogger(logger))
```

## 2. Métodos HTTP
O Goxios suporta todos os métodos HTTP comuns de forma fluida. ([Exemplos](cmd/examples/methods))

```go
client.Get("/get").Do(ctx)
client.Post("/post", bodyBytes).Do(ctx)
client.Put("/put", bodyBytes).Do(ctx)
client.Patch("/patch", bodyBytes).Do(ctx)
client.Delete("/delete").Do(ctx)
client.Head("/head").Do(ctx)
client.Options("/options").Do(ctx)
```

## 3. Configuração da Requisição (Fluent API)
Você pode customizar cada requisição individualmente antes de executá-la com `.Do()`. ([Exemplos](cmd/examples/headers/request))

```go
resp, err := client.Post("/resource", data).
    Header("X-Custom-Req", "true").
    Headers(map[string]string{"X-Another": "value"}).
    Do(ctx)
```

## 4. Autenticação

### Basic Auth
([Exemplo](cmd/examples/auth/basicauth))
```go
client, _ := goxios.New(goxios.WithBasicAuth("usuario", "senha"))
```

### Bearer Token
```go
client, _ := goxios.New(goxios.WithBearerToken("seu-token-aqui"))
```

### OAuth2 Client Credentials
Suporta renovação automática de tokens e cache. ([Exemplo](cmd/examples/auth/oauth/basic))
```go
client, _ := goxios.New(
    goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
        TokenURL:     "https://auth.server.com/token",
        ClientID:     "id",
        ClientSecret: "secret",
    }),
)
```

#### OAuth2 com Cache Redis
Permite compartilhar o token entre múltiplas instâncias da aplicação. ([Exemplo](cmd/examples/auth/oauth/cache_redis))
```go
import "github.com/drummerzzz/goxios/src/cache/redis"

cache := redis.NewRedisCache("localhost:6379")
goxios.WithOAuthClientCredentials(goxios.OAuthClientCredentialsConfig{
    // ...
    Cache: cache,
})
```

## 5. Trabalhando com Respostas
A struct `Response` oferece métodos para facilitar o consumo dos dados. ([Exemplos](cmd/examples/responses))

### Verificação de Sucesso
```go
if resp.Ok() {
    // Status code entre 200-299
}
```

### JSON para Struct (Generics)
```go
type User struct { ID int; Name string }
user, err := goxios.JsonAs[User](resp)
```

### JSON para Map
```go
data, err := resp.JsonMap()
fmt.Println(data["name"])
```

### Bytes Brutos
```go
body, err := resp.RawBytes()
```

## 6. mTLS (Mutual TLS)
O Goxios simplifica o uso de certificados digitais. ([Exemplos](cmd/examples/auth/mtls))

### Globalmente no Cliente
```go
// Via arquivos
goxios.WithMTLSFromFile("cert.pem", "key.pem")

// Via Base64
goxios.WithMTLSFromBase64("base64_cert", "base64_key")
```

### Por Requisição
```go
client.Get("/private").
    MTLS(&goxios.Certificate{
        MtlsCertBase64: "...",
        MtlsKeyBase64:  "...",
    }).
    Do(ctx)
```

