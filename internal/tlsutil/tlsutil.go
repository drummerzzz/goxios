package tlsutil

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
)

// LoadTLSConfigFromBase64 carrega uma configuração TLS a partir de certificados codificados em Base64.
func LoadTLSConfigFromBase64(certBase64, keyBase64 string) (*tls.Config, error) {
	if certBase64 == "" || keyBase64 == "" {
		return nil, errors.New("empty base64 cert/key")
	}

	certPEM, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return nil, err
	}
	keyPEM, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

