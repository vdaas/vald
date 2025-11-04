package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// helper to create a self-signed certificate/key pair and write to temp files.
func createTempCert(
	t *testing.T, dir string, serial int64,
) (certPath, keyPath string, cert *x509.Certificate, key *rsa.PrivateKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject:      pkix.Name{CommonName: "unit-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IsCA:         true,
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}

	certFile, err := os.CreateTemp(dir, "cert*.pem")
	if err != nil {
		t.Fatalf("temp cert: %v", err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	certFile.Close()

	keyFile, err := os.CreateTemp(dir, "key*.pem")
	if err != nil {
		t.Fatalf("temp key: %v", err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	keyFile.Close()

	parsedCert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	return certFile.Name(), keyFile.Name(), parsedCert, key
}

// createCRL writes a CRL containing revoked serials signed by issuer.
func createCRL(
	t *testing.T,
	dir string,
	issuer *x509.Certificate,
	issuerKey *rsa.PrivateKey,
	revokedSerials ...*big.Int,
) string {
	t.Helper()

	revoked := make([]pkix.RevokedCertificate, 0, len(revokedSerials))
	for _, sn := range revokedSerials {
		revoked = append(revoked, pkix.RevokedCertificate{
			SerialNumber:   sn,
			RevocationTime: time.Now().Add(-time.Minute),
		})
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		SignatureAlgorithm:  issuer.SignatureAlgorithm,
		Issuer:              issuer.Subject,
		RevokedCertificates: revoked,
		Number:              big.NewInt(1),
		ThisUpdate:          time.Now().Add(-time.Minute),
		NextUpdate:          time.Now().Add(time.Hour),
	}, issuer, issuerKey)
	if err != nil {
		t.Fatalf("create CRL: %v", err)
	}

	path := filepath.Join(dir, "revoked.crl")
	f, _ := os.Create(path)
	pem.Encode(f, &pem.Block{Type: "X509 CRL", Bytes: crlBytes})
	f.Close()
	return path
}

// TestCRLRevocation ensures that a revoked certificate is rejected by NewServerConfig.
func TestCRLRevocation(t *testing.T) {
	dir := t.TempDir()

	certPath, keyPath, cert, key := createTempCert(t, dir, 1001)

	crlPath := createCRL(t, dir, cert, key, big.NewInt(1001))

	opts := []Option{
		WithCert(certPath),
		WithKey(keyPath),
		WithCRL(crlPath),
	}

	if _, err := NewServerConfig(opts...); err == nil || !errors.Is(err, ErrCertRevoked) {
		t.Fatalf("expected ErrCertRevoked, got %v", err)
	}
}

// TestCRLSuccess ensures that non-revoked cert passes.
func TestCRLSuccess(t *testing.T) {
	dir := t.TempDir()

	certPath, keyPath, cert, key := createTempCert(t, dir, 2001)

	// CRL contains different serial -> should pass
	crlPath := createCRL(t, dir, cert, key, big.NewInt(9999))

	opts := []Option{
		WithCert(certPath),
		WithKey(keyPath),
		WithCRL(crlPath),
	}

	if _, err := NewServerConfig(opts...); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}
