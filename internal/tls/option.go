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

import "crypto/tls"

type Option func(*credentials) error

var (
	defaultOpts = []Option{
	}
)

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
		c.cfg = cfg
		return nil
	}
}
