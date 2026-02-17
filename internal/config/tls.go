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

package config

import "github.com/vdaas/vald/internal/tls"

// TLS represent the TLS configuration for server.
type TLS struct {
	// Cert represents the certificate file path.
	Cert string `json:"cert" yaml:"cert"`
	// Key represents the key file path.
	Key string `json:"key" yaml:"key"`
	// CA represents the CA certificate file path.
	CA string `json:"ca" yaml:"ca"`
	// CRL represents the CRL file path.
	CRL string `json:"crl" yaml:"crl"`
	// ServerName represents the server name.
	ServerName string `json:"server_name" yaml:"server_name"`
	// ClientAuth represents the client authentication type (NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven, RequireAndVerifyClientCert).
	ClientAuth string `json:"client_auth" yaml:"client_auth"`
	// Enabled enables TLS.
	Enabled bool `json:"enabled" yaml:"enabled"`
	// InsecureSkipVerify enables skipping verification.
	InsecureSkipVerify bool `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`
	// HotReload enables hot reload of certificates.
	HotReload bool `json:"hot_reload" yaml:"hot_reload"`
}

// Bind returns TLS object whose every value except Enabled is field value of environment value.
func (t *TLS) Bind() *TLS {
	t.Cert = GetActualValue(t.Cert)
	t.Key = GetActualValue(t.Key)
	t.CA = GetActualValue(t.CA)
	t.CRL = GetActualValue(t.CRL)
	t.ServerName = GetActualValue(t.ServerName)
	t.ClientAuth = GetActualValue(t.ClientAuth)
	return t
}

// Opts returns []tls.Option object whose every value is field value.
func (t *TLS) Opts() []tls.Option {
	if !t.Enabled {
		return nil
	}
	t = t.Bind()
	return []tls.Option{
		tls.WithCa(t.CA),
		tls.WithCert(t.Cert),
		tls.WithKey(t.Key),
		tls.WithCRL(t.CRL),
		tls.WithInsecureSkipVerify(t.InsecureSkipVerify),
		tls.WithServerName(t.ServerName),
		tls.WithClientAuth(t.ClientAuth),
		tls.WithHotReload(t.HotReload),
	}
}
