package goxios

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/drummerzzz/goxios/internal/tlsutil"
	goxios_errors "github.com/drummerzzz/goxios/src/errors"
	"github.com/drummerzzz/goxios/src/request"
	"github.com/drummerzzz/goxios/src/response"
	"go.uber.org/zap"
)

// Request alias para facilitar o uso sem import direto do pacote request
type Request = request.Request

// Response alias para facilitar o uso sem import direto do pacote response
type Response = response.Response

// JsonAs re-exporta a função JsonAs do pacote response.
func JsonAs[T any](r *Response) (T, error) {
	return response.JsonAs[T](r)
}

func (c *Client) Request(method, rawURL string) *Request {
	h := make(http.Header, len(c.defaultHeaders))
	for k, v := range c.defaultHeaders {
		h[k] = append([]string(nil), v...)
	}
	return &Request{
		HTTPClient:     c.httpClient,
		Transport:      c.transport,
		Logger:         c.logger,
		BaseURL:        c.baseURL,
		Method:         method,
		RawURL:         rawURL,
		CustomHeaders:  h,
		Auth:           c.defaultAuth,
		ErrNilClient:   goxios_errors.ErrNilClient,
		ErrEmptyURL:    goxios_errors.ErrEmptyURL,
		ErrRelativeURL: goxios_errors.ErrRelativeURL,
	}
}

func (c *Client) Get(rawURL string) *Request {
	return c.Request(http.MethodGet, rawURL)
}

func (c *Client) Post(rawURL string, body []byte) *Request {
	return c.Request(http.MethodPost, rawURL).Body(body)
}

func (c *Client) Put(rawURL string, body []byte) *Request {
	return c.Request(http.MethodPut, rawURL).Body(body)
}

func (c *Client) Patch(rawURL string, body []byte) *Request {
	return c.Request(http.MethodPatch, rawURL).Body(body)
}

func (c *Client) Delete(rawURL string) *Request {
	return c.Request(http.MethodDelete, rawURL)
}

func (c *Client) Head(rawURL string) *Request {
	return c.Request(http.MethodHead, rawURL)
}

func (c *Client) Options(rawURL string) *Request {
	return c.Request(http.MethodOptions, rawURL)
}

type Client struct {
	baseURL        *url.URL
	httpClient     *http.Client
	transport      *http.Transport
	defaultHeaders http.Header
	defaultAuth    request.AuthFunc
	logger         *zap.Logger
}

type Option func(*Client) error

func New(opts ...Option) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: nil,
	}

	c := &Client{
		httpClient: &http.Client{
			Transport: tr,
		},
		transport:      tr,
		defaultHeaders: make(http.Header),
		logger:         zap.NewNop(),
	}

	c.defaultHeaders.Set("Accept", "application/json")
	c.defaultHeaders.Set("Content-Type", "application/json")

	applyProxyFromEnv(c)

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithLogger define o logger padrão do client.
// Se logger for nil, usa zap.NewNop().
func WithLogger(logger *zap.Logger) Option {
	return func(c *Client) error {
		if c == nil {
			return nil
		}
		if logger == nil {
			c.logger = zap.NewNop()
			return nil
		}
		c.logger = logger
		return nil
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d < 0 {
			return goxios_errors.ErrInvalidTimeout
		}
		c.httpClient.Timeout = d
		return nil
	}
}

func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		if raw == "" {
			c.baseURL = nil
			return nil
		}
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		if u.Scheme == "" || u.Host == "" {
			return goxios_errors.ErrInvalidBaseURL
		}
		c.baseURL = u
		return nil
	}
}

func WithDefaultHeader(key, value string) Option {
	return func(c *Client) error {
		if key == "" {
			return goxios_errors.ErrEmptyHeaderKey
		}
		c.defaultHeaders.Set(key, value)
		return nil
	}
}

func WithDefaultHeaders(headers map[string]string) Option {
	return func(c *Client) error {
		for k, v := range headers {
			if k == "" {
				return goxios_errors.ErrEmptyHeaderKey
			}
			c.defaultHeaders.Set(k, v)
		}
		return nil
	}
}

func WithProxyURL(raw string) Option {
	return func(c *Client) error {
		if raw == "" {
			c.transport.Proxy = nil
			return nil
		}
		parsed, err := url.Parse(raw)
		if err != nil {
			return err
		}
		c.transport.Proxy = http.ProxyURL(parsed)

		if c.transport.TLSClientConfig == nil {
			c.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
		}
		return nil
	}
}

func WithMTLSFromBase64(certBase64, keyBase64 string) Option {
	return func(c *Client) error {
		if certBase64 == "" || keyBase64 == "" {
			c.transport.TLSClientConfig = nil
			return nil
		}
		tlsConfig, err := tlsutil.LoadTLSConfigFromBase64(certBase64, keyBase64)
		if err != nil {
			return err
		}
		c.transport.TLSClientConfig = tlsConfig
		return nil
	}
}

func WithMTLSFromFile(certFile, keyFile string) Option {
	return func(c *Client) error {
		if certFile == "" || keyFile == "" {
			c.transport.TLSClientConfig = nil
			return nil
		}
		certBytes, err := os.ReadFile(certFile)
		if err != nil {
			return err
		}
		keyBytes, err := os.ReadFile(keyFile)
		if err != nil {
			return err
		}
		certBase64 := base64.StdEncoding.EncodeToString(certBytes)
		keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)
		return WithMTLSFromBase64(certBase64, keyBase64)(c)
	}
}

func applyProxyFromEnv(c *Client) {
	proxyURL := os.Getenv("GOXIOS_HTTP_PROXY")
	useProxy := os.Getenv("GOXIOS_USE_HTTP_PROXY") == "true"
	if !useProxy || proxyURL == "" {
		return
	}
	_ = WithProxyURL(proxyURL)(c)
}
