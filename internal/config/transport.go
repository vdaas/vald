//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// TCP represents the TCP configuration for server.
type Transport struct {
	RoundTripper *RoundTripper `yaml:"round_tripper" json:"round_tripper"`
	Backoff      *Backoff      `yaml:"backoff"       json:"backoff"`
}

// RoundTripper represents the round trip configuration for transport.
type RoundTripper struct {
	TLSHandshakeTimeout   string `yaml:"tls_handshake_timeout"    json:"tls_handshake_timeout"`
	MaxIdleConns          int    `yaml:"max_idle_conns"           json:"max_idle_conns"`
	MaxIdleConnsPerHost   int    `yaml:"max_idle_conns_per_host"  json:"max_idle_conns_per_host"`
	MaxConnsPerHost       int    `yaml:"max_conns_per_host"       json:"max_conns_per_host"`
	IdleConnTimeout       string `yaml:"idle_conn_timeout"        json:"idle_conn_timeout"`
	ResponseHeaderTimeout string `yaml:"response_header_timeout"  json:"response_header_timeout"`
	ExpectContinueTimeout string `yaml:"expect_continue_timeout"  json:"expect_continue_timeout"`
	MaxResponseHeaderSize int64  `yaml:"max_response_header_size" json:"max_response_header_size"`
	WriteBufferSize       int64  `yaml:"write_buffer_size"        json:"write_buffer_size"`
	ReadBufferSize        int64  `yaml:"read_buffer_size"         json:"read_buffer_size"`
	ForceAttemptHTTP2     bool   `yaml:"force_attempt_http_2"     json:"force_attempt_http_2"`
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
