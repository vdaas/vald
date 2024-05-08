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

import "crypto/tls"

type Option func(*credentials) error

var defaultOptions = func() []Option {
	return []Option{
		WithInsecureSkipVerify(false),
		WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
			NextProtos: []string{
				"http/1.1",
				"h2",
			},
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
				tls.X25519,
			},
			SessionTicketsDisabled: true,
			// PreferServerCipherSuites: true,
			// CipherSuites: []uint16{
			// tls.TLS_RSA_WITH_RC4_128_SHA,
			// tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			// tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, // Maybe this is work on TLS 1.2
			// tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA, // TLS1.3 Feature
			// tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA, // TLS1.3 Feature
			// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			// },
			ClientAuth: tls.NoClientCert,
		}),
	}
}

func WithCert(cert string) Option {
	return func(c *credentials) error {
		c.cert = cert
		return nil
	}
}

func WithKey(key string) Option {
	return func(c *credentials) error {
		c.key = key
		return nil
	}
}

func WithCa(ca string) Option {
	return func(c *credentials) error {
		c.ca = ca
		return nil
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(c *credentials) error {
		if cfg != nil {
			c.cfg = cfg
		}
		return nil
	}
}

func WithInsecureSkipVerify(insecure bool) Option {
	return func(c *credentials) error {
		c.insecure = insecure
		return nil
	}
}
