package request

import (
	"errors"
	"net/url"
	"testing"
)

func TestRequest_ResolveURL(t *testing.T) {
	baseURL, _ := url.Parse("https://api.example.com/v1/")
	r := &Request{
		BaseURL:        baseURL,
		ErrEmptyURL:    errors.New("empty url"),
		ErrRelativeURL: errors.New("relative url"),
	}

	tests := []struct {
		raw      string
		expected string
		wantErr  bool
	}{
		{"/users", "https://api.example.com/users", false},
		{"users", "https://api.example.com/v1/users", false},
		{"https://other.com/api", "https://other.com/api", false},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := r.ResolveURL(tt.raw)
		if (err != nil) != tt.wantErr {
			t.Errorf("ResolveURL(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("ResolveURL(%q) = %q, want %q", tt.raw, got, tt.expected)
		}
	}
}

