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

// TCP represents the TCP configuration for server.
type Transport struct {
	// RoundTripper represents the round tripper configuration.
	RoundTripper *RoundTripper `json:"round_tripper" yaml:"round_tripper"`
	// Backoff represents the backoff configuration.
	Backoff *Backoff `json:"backoff" yaml:"backoff"`
}

// RoundTripper represents the round trip configuration for transport.
type RoundTripper struct {
	// TLSHandshakeTimeout represents the TLS handshake timeout.
	TLSHandshakeTimeout string `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`
	// IdleConnTimeout represents the idle connection timeout.
	IdleConnTimeout string `json:"idle_conn_timeout" yaml:"idle_conn_timeout"`
	// ResponseHeaderTimeout represents the response header timeout.
	ResponseHeaderTimeout string `json:"response_header_timeout" yaml:"response_header_timeout"`
	// ExpectContinueTimeout represents the expect continue timeout.
	ExpectContinueTimeout string `json:"expect_continue_timeout" yaml:"expect_continue_timeout"`
	// MaxIdleConns represents the maximum idle connections.
	MaxIdleConns int `json:"max_idle_conns" yaml:"max_idle_conns"`
	// MaxIdleConnsPerHost represents the maximum idle connections per host.
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	// MaxConnsPerHost represents the maximum connections per host.
	MaxConnsPerHost int `json:"max_conns_per_host" yaml:"max_conns_per_host"`
	// MaxResponseHeaderSize represents the maximum response header size.
	MaxResponseHeaderSize int64 `json:"max_response_header_size" yaml:"max_response_header_size"`
	// WriteBufferSize represents the write buffer size.
	WriteBufferSize int64 `json:"write_buffer_size" yaml:"write_buffer_size"`
	// ReadBufferSize represents the read buffer size.
	ReadBufferSize int64 `json:"read_buffer_size" yaml:"read_buffer_size"`
	// ForceAttemptHTTP2 enables force attempt HTTP2.
	ForceAttemptHTTP2 bool `json:"force_attempt_http_2" yaml:"force_attempt_http_2"`
}

// Bind binds the actual data from the RoundTripper receiver fields.
func (r *RoundTripper) Bind() *RoundTripper {
	r.TLSHandshakeTimeout = GetActualValue(r.TLSHandshakeTimeout)
	r.IdleConnTimeout = GetActualValue(r.IdleConnTimeout)
	r.ResponseHeaderTimeout = GetActualValue(r.ResponseHeaderTimeout)
	r.ExpectContinueTimeout = GetActualValue(r.ExpectContinueTimeout)
	return r
}

// Bind binds the actual data from the Transport receiver fields.
func (t *Transport) Bind() *Transport {
	if t.RoundTripper != nil {
		t.RoundTripper = t.RoundTripper.Bind()
	} else {
		t.RoundTripper = new(RoundTripper)
	}

	if t.Backoff != nil {
		t.Backoff = t.Backoff.Bind()
	} else {
		t.Backoff = new(Backoff)
	}

	return t
}
