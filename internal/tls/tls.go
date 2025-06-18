//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package tls

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
)

type (
	// Conn alias for tls.Conn.
	Conn = tls.Conn
	// Config alias for tls.Config.
	Config = tls.Config
	// Dialer alias for tls.Dialer.
	Dialer = tls.Dialer
	// Certificate alias for tls.Certificate.
	Certificate = tls.Certificate
)

const (
	VersionTLS10 = tls.VersionTLS10
	VersionTLS11 = tls.VersionTLS11
	VersionTLS12 = tls.VersionTLS12
	VersionTLS13 = tls.VersionTLS13
)

var (
	Client         = tls.Client
	Dial           = tls.Dial
	DialWithDialer = tls.DialWithDialer
	NewListener    = tls.NewListener
	Listen         = tls.Listen
)

// credentials holds TLS settings for server and client
// including certificate paths, CA bundle, and hot reload policies.
type credentials struct {
	cfg        *Config
	cert       string
	key        string
	ca         string
	sn         string
	insecure   bool
	clientAuth tls.ClientAuthType
}

// newCredential builds credentials from defaults and provided options.
func newCredential(opts ...Option) (*credentials, error) {
	c := new(credentials)
	for _, opt := range append(defaultOptions(), opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.cfg == nil {
		c.cfg = new(Config)
	}
	if c.sn != "" {
		c.cfg.ServerName = c.sn
	}
	c.cfg.InsecureSkipVerify = c.insecure
	if c.cfg.MinVersion == 0 {
		c.cfg.MinVersion = tls.VersionTLS12
	}
	if c.cfg.CurvePreferences == nil {
		c.cfg.CurvePreferences = defaultCurvePreferences
	}
	if c.cfg.NextProtos == nil {
		c.cfg.NextProtos = defaultNextProtos
	}
	c.cfg.SessionTicketsDisabled = true

	if c.clientAuth == tls.NoClientCert {
		switch {
		case c.ca == "" && c.cert == "" && c.key == "":
			c.clientAuth = tls.NoClientCert
		case c.ca != "" && (c.cert == "" || c.key == ""):
			c.clientAuth = tls.VerifyClientCertIfGiven
		case c.ca != "" && c.cert != "" && c.key != "":
			c.clientAuth = tls.RequireAndVerifyClientCert
		default:
			c.clientAuth = tls.RequireAnyClientCert
		}
	}
	return c, nil
}

// loadKeyPair loads a certificate and key, wrapping errors.Wrapf.
func loadKeyPair(role, certPath, keyPath string) (tls.Certificate, error) {
	kp, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Warn(errors.Join(err, errors.ErrFailedToLoadCertKey(role, certPath, keyPath)))
		return tls.Certificate{}, err
	}
	return kp, nil
}

// NewServerConfig returns a *tls.Config for server, with optional mTLS and hot reload
func NewServerConfig(opts ...Option) (*Config, error) {
	c, err := newCredential(opts...)
	if err != nil {
		return nil, err
	}
	// require cert and key
	if c.cert == "" || c.key == "" {
		if c.insecure {
			return &Config{
				InsecureSkipVerify: c.insecure,
			}, nil
		}
		return nil, errors.ErrTLSCertOrKeyNotFound
	}
	if c.sn == "" {
		c.sn = "vald-server"
	}
	// load cert pair
	kp, err := loadKeyPair(c.sn, c.cert, c.key)
	if err != nil {
		return nil, err
	}
	c.cfg.Certificates = []tls.Certificate{kp}
	// if CA provided, configure mTLS
	if c.ca != "" {
		pool, err := NewX509CertPool(c.ca)
		if err != nil {
			return nil, err
		}
		c.cfg.ClientCAs = pool
		c.cfg.ClientAuth = c.clientAuth
	}
	return c.cfg, nil
}

// NewClientConfig returns a *tls.Config for client, with optional mTLS and hot reload
func NewClientConfig(opts ...Option) (*Config, error) {
	c, err := newCredential(opts...)
	if err != nil {
		return nil, err
	}

	if c.cert == "" && c.key == "" && c.ca == "" {
		if c.insecure {
			return &Config{
				InsecureSkipVerify: c.insecure,
			}, nil
		}
		return nil, errors.ErrTLSCertOrKeyNotFound
	}
	// setup RootCAs from CA bundle or self-signed cert
	if c.ca != "" || c.cert != "" {
		pool, err := NewX509CertPool(c.ca, c.cert)
		if err != nil {
			// Only return error if CA was explicitly provided
			if c.ca != "" {
				return nil, err
			}
			// Log the error when only cert was provided
			log.Warnf("Failed to create RootCAs pool from cert: %v", err)
		} else if pool != nil {
			c.cfg.RootCAs = pool
		}
	}
	// load client cert if provided
	if c.cert != "" && c.key != "" {
		if c.sn == "" {
			c.sn = "vald-client"
		}
		kp, err := loadKeyPair(c.sn, c.cert, c.key)
		if err != nil {
			return nil, err
		}
		c.cfg.Certificates = []tls.Certificate{kp}
	}
	return c.cfg, nil
}

// NewX509CertPool loads one or more PEM files into a CertPool
// It deduplicates certificates, logs SANs, checks expiration and chain
func NewX509CertPool(paths ...string) (*x509.CertPool, error) {
	pool := systemOrNewPool()
	seen := make(map[string]struct{}, len(paths))
	added := false

	for _, path := range paths {
		if path == "" {
			continue
		}
		data, err := file.ReadFile(path)
		if err != nil {
			log.Warnf("failed to read %s: %v", path, err)
			continue
		}
		certs, err := parsePEMCertificates(data)
		if err != nil {
			log.Warnf("failed to parse certificates in %s: %v", path, err)
			continue
		}
		for _, cert := range certs {
			if cert == nil {
				continue
			}
			if processCert(path, cert, pool, seen) {
				added = true
			}
		}
	}

	if !added {
		return nil, errors.ErrNoCertsAddedToPool
	}
	return pool, nil
}

// systemOrNewPool returns the system pool, or a new one if unavailable.
func systemOrNewPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		return x509.NewCertPool()
	}
	return pool
}

// processCert logs, verifies, and adds cert to pool if CA or self-signed, deduplicating by fingerprint.
func processCert(
	path string, cert *x509.Certificate, pool *x509.CertPool, seen map[string]struct{},
) bool {
	now := time.Now()
	checkSignature := cert.CheckSignatureFrom(cert)
	selfSigned := checkSignature == nil
	log.Debugf("Loaded cert Path=%s CN=%s Issuer=%s IsCA=%t SelfSigned=%t Expired=%t(%s) SANs(DNS: %v, IP: %v), SignatureCheckResult=\"%v\"",
		path,
		cert.Subject.CommonName,
		cert.Issuer.CommonName,
		cert.IsCA,
		selfSigned,
		now.After(cert.NotAfter),
		cert.NotAfter,
		cert.DNSNames,
		cert.IPAddresses,
		checkSignature)

	if err := verifyCertChain(cert, pool, now); err != nil {
		log.Warnf("chain verify failed for %s: %v", cert.Subject.CommonName, err)
	}

	if !cert.IsCA && !selfSigned {
		return false
	}
	fp := fingerprint(cert.Raw)
	if _, exists := seen[fp]; exists {
		return false
	}
	pool.AddCert(cert)
	seen[fp] = struct{}{}
	return true
}

// verifyCertChain attempts to verify the cert against the provided pool.
func verifyCertChain(cert *x509.Certificate, pool *x509.CertPool, now time.Time) error {
	opts := x509.VerifyOptions{
		Roots:         pool,
		Intermediates: x509.NewCertPool(),
		CurrentTime:   now,
	}
	_, err := cert.Verify(opts)
	return err
}

// fingerprint returns the SHA-256 hex of data.
func fingerprint(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// parsePEMCertificates extracts x509.Certificates from PEM data,
// skipping private keys and unknown blocks
func parsePEMCertificates(pemBytes []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			break
		}
		if len(block.Headers) != 0 {
			continue
		}
		switch block.Type {
		case "CERTIFICATE":
			cs, err := x509.ParseCertificates(block.Bytes)
			if err != nil {
				log.Warn(errors.Wrap(err, "x509.ParseCertificate in parsePEMCertificates returned error"))
				continue
			}
			certs = append(certs, cs...)
		case "PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY", "ED25519 PRIVATE KEY":
			log.Debugf("Skipping key block: %s", block.Type)
		default:
			log.Debugf("Unknown PEM block: %s", block.Type)
		}
	}
	if len(certs) == 0 {
		return nil, errors.ErrNoCertsFoundInPEM
	}
	return certs, nil
}
