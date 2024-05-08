//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package tls provides implementation of Go API for tls certificate provider
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
)

type (
	Conn   = tls.Conn
	Config = tls.Config
	Dialer = tls.Dialer
)

type credentials struct {
	cfg      *tls.Config
	cert     string
	key      string
	ca       string
	insecure bool
}

var (
	Client         = tls.Client
	Dial           = tls.Dial
	DialWithDialer = tls.DialWithDialer
)

// NewTLSConfig returns a *tls.Config struct or error
// This function read TLS configuration and initialize *tls.Config struct.
// This function initialize TLS configuration, for example the CA certificate and key to start TLS server.
// Server and CA Certificate, and private key will read from a file from the file path definied in environment variable.
func New(opts ...Option) (*Config, error) {
	c, err := newCredential(opts...)
	if err != nil {
		return nil, err
	}

	if c.cert == "" || c.key == "" {
		return nil, errors.ErrTLSCertOrKeyNotFound
	}

	c.cfg.Certificates = make([]tls.Certificate, 1)
	c.cfg.Certificates[0], err = tls.LoadX509KeyPair(c.cert, c.key)
	if err != nil {
		return nil, err
	}

	if c.ca != "" {
		pool, err := NewX509CertPool(c.ca)
		if err != nil {
			return nil, err
		}
		c.cfg.ClientCAs = pool
		c.cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return c.cfg, nil
}

func NewClientConfig(opts ...Option) (*Config, error) {
	c, err := newCredential(opts...)
	if err != nil {
		return nil, err
	}

	if c.ca != "" {
		pool, err := NewX509CertPool(c.ca)
		if err != nil {
			return nil, err
		}
		c.cfg.RootCAs = pool
	}

	if c.cert != "" && c.key != "" {
		c.cfg.Certificates = make([]tls.Certificate, 1)
		c.cfg.Certificates[0], err = tls.LoadX509KeyPair(c.cert, c.key)
		if err != nil {
			return nil, err
		}
	}

	return c.cfg, nil
}

// NewX509CertPool returns *x509.CertPool struct or error.
// The CertPool will read the certificate from the path, and append the content to the system certificate pool, and return.
func NewX509CertPool(path string) (pool *x509.CertPool, err error) {
	c, err := file.ReadFile(path)
	if err != nil || c == nil {
		return nil, err
	}
	if err == nil && c != nil {
		pool, err = x509.SystemCertPool()
		if err != nil || pool == nil {
			pool = x509.NewCertPool()
		}
		if !pool.AppendCertsFromPEM(c) {
			err = errors.ErrCertificationFailed
		}
	}
	return pool, err
}

func newCredential(opts ...Option) (c *credentials, err error) {
	c = new(credentials)

	for _, opt := range append(defaultOptions(), opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	c.cfg.InsecureSkipVerify = c.insecure
	return c, nil
}
