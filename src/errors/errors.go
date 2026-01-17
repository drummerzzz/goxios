package goxios_errors

import "errors"

var (
	ErrInvalidBaseURL  = errors.New("baseURL must be absolute (e.g. https://api.example.com)")
	ErrEmptyURL        = errors.New("empty url")
	ErrRelativeURL     = errors.New("relative url without baseURL configured")
	ErrNilClient       = errors.New("request/client is nil")
	ErrInvalidTimeout  = errors.New("invalid timeout")
	ErrEmptyHeaderKey  = errors.New("empty header key")
	ErrUnsupportedAuth = errors.New("unsupported auth type")
)
