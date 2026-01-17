package tlsutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func generateTestCert() (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return base64.StdEncoding.EncodeToString(certPem), base64.StdEncoding.EncodeToString(keyPem), nil
}

func TestLoadTLSConfigFromBase64(t *testing.T) {
	t.Run("ValidCertAndKey", func(t *testing.T) {
		certBase64, keyBase64, err := generateTestCert()
		if err != nil {
			t.Fatalf("failed to generate test cert: %v", err)
		}

		cfg, err := LoadTLSConfigFromBase64(certBase64, keyBase64)
		if err != nil {
			t.Fatalf("LoadTLSConfigFromBase64 failed with error: %v", err)
		}

		if cfg == nil {
			t.Fatal("expected non-nil config")
		}

		if len(cfg.Certificates) != 1 {
			t.Errorf("expected 1 certificate, got %d", len(cfg.Certificates))
		}
	})

	t.Run("EmptyInputs", func(t *testing.T) {
		_, err := LoadTLSConfigFromBase64("", "")
		if err == nil {
			t.Error("expected error for empty inputs")
		}
	})

	t.Run("InvalidBase64", func(t *testing.T) {
		_, err := LoadTLSConfigFromBase64("!!!", "!!!")
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})

	t.Run("InvalidPEM", func(t *testing.T) {
		invalidCert := base64.StdEncoding.EncodeToString([]byte("not a certificate"))
		invalidKey := base64.StdEncoding.EncodeToString([]byte("not a key"))
		_, err := LoadTLSConfigFromBase64(invalidCert, invalidKey)
		if err == nil {
			t.Error("expected error for invalid PEM content")
		}
	})
}

