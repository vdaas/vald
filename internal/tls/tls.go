// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
