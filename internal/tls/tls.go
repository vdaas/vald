//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"os"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
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

	for _, opt := range append(defaultOptions(), opts...) {
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

	return c.cfg, nil
}

func NewClientConfig(opts ...Option) (*Config, error) {
	var err error
	c := new(credentials)

	for _, opt := range append(defaultOptions(), opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
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
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var n int64 = bytes.MinRead
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, n))
	err = safety.RecoverFunc(func() (err error) {
		_, err = buf.ReadFrom(f)
		return err
	})()
	if err != nil {
		return nil, err
	}
	c := buf.Bytes()
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
