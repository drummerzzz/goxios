package goxios

import (
	"testing"
	"time"
)

func TestClient_Options(t *testing.T) {
	baseURL := "https://api.example.com"
	timeout := 5 * time.Second

	c, err := New(
		WithBaseURL(baseURL),
		WithTimeout(timeout),
		WithDefaultHeader("X-Test", "Value"),
	)
	if err != nil {
		t.Fatalf("New() errored: %v", err)
	}

	if c.baseURL.String() != baseURL {
		t.Errorf("expected baseURL %s, got %s", baseURL, c.baseURL.String())
	}
	if c.httpClient.Timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, c.httpClient.Timeout)
	}
	if c.defaultHeaders.Get("X-Test") != "Value" {
		t.Errorf("expected header X-Test Value, got %s", c.defaultHeaders.Get("X-Test"))
	}
}

func TestClient_WithMTLS_Nil(t *testing.T) {
	c, _ := New(WithMTLSFromBase64("", ""))
	if c.transport.TLSClientConfig != nil {
		t.Error("expected TLSClientConfig nil when passing nil cert")
	}
}
