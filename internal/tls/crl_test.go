//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package tls_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// helper to create a self-signed certificate/key pair and write to temp files.
func createTempCert(t *testing.T) (cert *x509.Certificate, key *rsa.PrivateKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(9999),
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

	parsedCert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}

	return parsedCert, key
}

// createCRL writes a CRL containing revoked serials signed by issuer.
func createCRL(t *testing.T, dir string, revokedSerials ...*big.Int) string {
	t.Helper()

	revoked := make([]pkix.RevokedCertificate, 0, len(revokedSerials))
	for _, sn := range revokedSerials {
		revoked = append(revoked, pkix.RevokedCertificate{
			SerialNumber:   sn,
			RevocationTime: time.Now().Add(-time.Minute),
		})
	}

	issuer, issuerKey := createTempCert(t)

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
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("os.Create: %v", err)
	}
	if err := pem.Encode(f, &pem.Block{Type: "X509 CRL", Bytes: crlBytes}); err != nil {
		f.Close()
		t.Fatalf("pem encode: %v", err)
	}
	f.Close()
	return path
}

func TestCRLRevocation(t *testing.T) {
	ctx, stop, addr := serverStarter(t, false)
	defer stop()

	dir := t.TempDir()
	serverCertPath := test.GetTestdataPath("tls/server.crt")
	data, err := os.ReadFile(serverCertPath)
	if err != nil {
		t.Fatalf("read server cert: %v", err)
	}
	block, _ := pem.Decode(data)
	srvCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse server cert: %v", err)
	}
	crlPath := createCRL(t, dir, srvCert.SerialNumber)

	ccfg, err := tls.NewClientConfig(
		tls.WithCa(test.GetTestdataPath("tls/ca.pem")),
		tls.WithCRL(crlPath),
	)
	if err != nil {
		t.Fatalf("client tls: %v", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(ccfg)))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	_, err = vald.NewIndexClient(conn).IndexInfo(ctx, &payload.Empty{})
	if !strings.Contains(err.Error(), "certificate revoked") {
		t.Fatalf("expected ErrCertRevoked, got %v", err)
	}
}

func TestCRLSuccess(t *testing.T) {
	ctx, stop, addr := serverStarter(t, false)
	defer stop()

	dir := t.TempDir()
	serverCertPath := test.GetTestdataPath("tls/server.crt")
	data, err := os.ReadFile(serverCertPath)
	if err != nil {
		t.Fatalf("read server cert: %v", err)
	}
	block, _ := pem.Decode(data)
	srvCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse server cert: %v", err)
	}
	otherSerial := new(big.Int).Add(srvCert.SerialNumber, big.NewInt(10))
	crlPath := createCRL(t, dir, otherSerial)

	ccfg, err := tls.NewClientConfig(
		tls.WithCa(test.GetTestdataPath("tls/ca.pem")),
		tls.WithCRL(crlPath),
	)
	if err != nil {
		t.Fatalf("client tls: %v", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(ccfg)))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	_, err = vald.NewIndexClient(conn).IndexInfo(ctx, &payload.Empty{})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}
