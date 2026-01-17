package goxios

import (
	"net/http"

	"github.com/drummerzzz/goxios/src/auth/oauth"
	"github.com/drummerzzz/goxios/src/request"
)

// OAuthClientCredentialsConfig re-exporta Config do pacote oauth para conveniência.
// Por padrão usa oauth.DefaultTokenResponse.
type OAuthClientCredentialsConfig = oauth.Config[oauth.DefaultTokenResponse]

// WithOAuthClientCredentials define auth default do client via OAuth2 client_credentials.
func WithOAuthClientCredentials(cfg OAuthClientCredentialsConfig) Option {
	return func(c *Client) error {
		if c == nil {
			return nil
		}
		src := oauth.NewTokenSource(c.httpClient, c.logger, cfg)
		c.defaultAuth = src.Apply
		return nil
	}
}

// WithOAuthClientCredentialsCustom permite definir auth OAuth2 com uma struct de resposta customizada.
func WithOAuthClientCredentialsCustom[T oauth.TokenResponse](cfg oauth.Config[T]) Option {
	return func(c *Client) error {
		if c == nil {
			return nil
		}
		src := oauth.NewTokenSource(c.httpClient, c.logger, cfg)
		c.defaultAuth = src.Apply
		return nil
	}
}

// Certificate re-exporta o tipo Certificate do pacote request para conveniência.
type Certificate = request.Certificate

func WithBasicAuth(username, password string) Option {
	return func(c *Client) error {
		if c == nil {
			return nil
		}
		c.defaultAuth = func(req *http.Request) error {
			if req == nil {
				return nil
			}
			req.SetBasicAuth(username, password)
			return nil
		}
		return nil
	}
}

func WithBearerToken(token string) Option {
	return func(c *Client) error {
		if c == nil {
			return nil
		}
		if token == "" {
			c.defaultAuth = nil
			return nil
		}
		c.defaultAuth = func(req *http.Request) error {
			if req == nil {
				return nil
			}
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}
		return nil
	}
}
