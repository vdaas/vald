// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
)

// crypto/tls calls VerifyPeerCertificate with no verified chains when client
// certificates are requested but not verified. The CRL hook used to reject all
// of those handshakes as revoked; it must only reject a revoked leaf.
func TestAttachCRLChainChecker_noVerifiedChains(t *testing.T) {
	t.Parallel()
	b, err := os.ReadFile("../test/data/tls/client.crt")
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(b)
	if block == nil {
		t.Fatal("no PEM block in client.crt")
	}
	leaf, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	hook := func(revoked map[string]struct{}) func([][]byte, [][]*x509.Certificate) error {
		c := &credentials{
			crl:           "unused.crl",
			revoked:       revoked,
			crlNextUpdate: time.Now().Add(time.Hour),
		}
		cfg := new(tls.Config)
		c.attachCRLChainChecker(cfg)
		return cfg.VerifyPeerCertificate
	}
	tests := []struct {
		revoked map[string]struct{}
		want    error
		name    string
		raw     [][]byte
	}{
		{name: "no client certificate", revoked: map[string]struct{}{"1": {}}, raw: nil, want: nil},
		{name: "unrevoked leaf", revoked: map[string]struct{}{"1": {}}, raw: [][]byte{block.Bytes}, want: nil},
		{name: "revoked leaf", revoked: map[string]struct{}{leaf.SerialNumber.String(): {}}, raw: [][]byte{block.Bytes}, want: errors.ErrCertRevoked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := hook(tt.revoked)(tt.raw, nil); !errors.Is(err, tt.want) {
				t.Errorf("VerifyPeerCertificate() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// ensureCRL used to drop the backoff that loadCRL returns on a read error, so with an
// expired NextUpdate and a missing CRL file every handshake reloaded under the write lock.
func TestEnsureCRL_readErrorBacksOff(t *testing.T) {
	t.Parallel()
	revoked := map[string]struct{}{"1": {}}
	c := &credentials{
		crl:           t.TempDir() + "/missing.crl",
		revoked:       revoked,
		crlNextUpdate: time.Now().Add(-time.Hour),
	}
	c.ensureCRL()
	if !c.crlNextUpdate.After(time.Now()) {
		t.Errorf("crlNextUpdate = %v, want a time in the future after a failed reload", c.crlNextUpdate)
	}
	if _, ok := c.revoked["1"]; !ok || len(c.revoked) != len(revoked) {
		t.Errorf("revoked = %v, want the previous set %v to be kept", c.revoked, revoked)
	}
}
