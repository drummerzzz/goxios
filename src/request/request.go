package request

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/drummerzzz/goxios/internal/tlsutil"
	"github.com/drummerzzz/goxios/src/response"
	"go.uber.org/zap"
)

// Certificate define os certificados para mTLS.
type Certificate struct {
	MtlsCertBase64 string
	MtlsKeyBase64  string
}

// AuthFunc aplica autenticação na request.
type AuthFunc func(req *http.Request) error

type Request struct {
	Method        string
	RawURL        string
	BodyData      []byte
	CustomHeaders http.Header
	Auth          AuthFunc
	MtlsCert      *Certificate
	HTTPClient    *http.Client
	Transport     *http.Transport
	Logger        *zap.Logger
	BaseURL       *url.URL

	// Erros pré-definidos para evitar ciclo de importação
	ErrNilClient   error
	ErrEmptyURL    error
	ErrRelativeURL error
}

func (r *Request) Header(key, value string) *Request {
	if r == nil {
		return r
	}
	if r.CustomHeaders == nil {
		r.CustomHeaders = make(http.Header)
	}
	r.CustomHeaders.Set(key, value)
	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	if r == nil {
		return r
	}
	for k, v := range headers {
		r.Header(k, v)
	}
	return r
}

func (r *Request) WithAuth(auth AuthFunc) *Request {
	if r == nil {
		return r
	}
	r.Auth = auth
	return r
}

// MTLS configura mTLS apenas para essa request.
func (r *Request) MTLS(cert *Certificate) *Request {
	if r == nil {
		return r
	}
	r.MtlsCert = cert
	return r
}

func (r *Request) Body(body []byte) *Request {
	if r == nil {
		return r
	}
	r.BodyData = body
	return r
}

func (r *Request) ResolveURL(raw string) (string, error) {
	if raw == "" {
		return "", r.ErrEmptyURL
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw, nil
	}
	if r.BaseURL == nil {
		return "", r.ErrRelativeURL
	}

	rel, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	return r.BaseURL.ResolveReference(rel).String(), nil
}

func (r *Request) Do(ctx ...context.Context) (*response.Response, error) {
	if r == nil || r.HTTPClient == nil {
		return nil, r.ErrNilClient
	}

	var c context.Context
	if len(ctx) > 0 {
		c = ctx[0]
	}
	if c == nil {
		c = context.Background()
	}

	finalURL, err := r.ResolveURL(r.RawURL)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if r.BodyData != nil {
		bodyReader = bytes.NewReader(r.BodyData)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(c, r.Method, finalURL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range r.CustomHeaders {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	if r.Auth != nil {
		if err := r.Auth(req); err != nil {
			if r.Logger != nil {
				r.Logger.Debug(
					"goxios request: error applying auth",
					zap.String("method", r.Method),
					zap.String("url", finalURL),
					zap.Error(err),
				)
			}
			return nil, err
		}
	}

	httpClient := r.HTTPClient
	if r.MtlsCert != nil {
		tr := r.Transport.Clone()
		tlsConfig, err := tlsutil.LoadTLSConfigFromBase64(r.MtlsCert.MtlsCertBase64, r.MtlsCert.MtlsKeyBase64)
		if err != nil {
			return nil, err
		}
		tr.TLSClientConfig = tlsConfig
		httpClient = &http.Client{
			Transport: tr,
			Timeout:   r.HTTPClient.Timeout,
		}
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		if r.Logger != nil {
			r.Logger.Debug(
				"goxios request: error executing request",
				zap.String("method", r.Method),
				zap.String("url", finalURL),
				zap.Duration("duration", time.Since(start)),
				zap.Error(err),
			)
		}
		return nil, err
	}

	return &response.Response{Response: resp, Logger: r.Logger}, nil
}
