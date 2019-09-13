//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"io/ioutil"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
)

type Config = tls.Config

type credentials struct {
	cfg  *tls.Config
	cert string
	key  string
	ca   string
}

// NewTLSConfig returns a *tls.Config struct or error
// This function read TLS configuration and initialize *tls.Config struct.
// This function initialize TLS configuration, for example the CA certificate and key to start TLS server.
// Server and CA Certificate, and private key will read from a file from the file path definied in environment variable.
func New(opts ...Option) (*Config, error) {
	var err error
	c := new(credentials)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
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

	c.cfg.BuildNameToCertificate()
	return c.cfg, nil
}

// NewX509CertPool returns *x509.CertPool struct or error.
// The CertPool will read the certificate from the path, and append the content to the system certificate pool, and return.
func NewX509CertPool(path string) (*x509.CertPool, error) {
	var pool *x509.CertPool
	c, err := ioutil.ReadFile(path)
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
