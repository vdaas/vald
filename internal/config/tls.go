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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/tls"

// TLS represent the TLS configuration for server.
type TLS struct {
	// Enable represent the server enable TLS or not.
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Cert represent the certificate environment variable key used to start server.
	Cert string `json:"cert" yaml:"cert"`

	// Key represent the private key environment variable key used to start server.
	Key string `json:"key" yaml:"key"`

	// CA represent the CA certificate environment variable key used to start server.
	CA string `json:"ca" yaml:"ca"`

	// InsecureSkipVerify represent enable/disable skip SSL certificate verification
	InsecureSkipVerify bool `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`
}

// Bind returns TLS object whose every value except Enabled is field value of environment value.
func (t *TLS) Bind() *TLS {
	t.Cert = GetActualValue(t.Cert)
	t.Key = GetActualValue(t.Key)
	t.CA = GetActualValue(t.CA)
	return t
}

// Opts returns []tls.Option object whose every value is field value.
func (t *TLS) Opts() []tls.Option {
	return []tls.Option{
		tls.WithCa(t.CA),
		tls.WithCert(t.Cert),
		tls.WithKey(t.Key),
		tls.WithInsecureSkipVerify(t.InsecureSkipVerify),
	}
}
