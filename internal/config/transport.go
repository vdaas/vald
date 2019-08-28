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

// Package config providers configuration type and load configuration logic
package config

// TCP represent the TCP configuration for server.
type Transport struct {
	RoundTripper struct {
		TLSHandshakeTimeout   string `yaml:"tls_handshake_timeout" json:"tls_handshake_timeout"`
		MaxIdleConns          int    `yaml:"max_idle_conns" json:"max_idle_conns"`
		MaxIdleConnsPerHost   int    `yaml:"max_idle_conns_per_host" json:"max_idle_conns_per_host"`
		MaxConnsPerHost       int    `yaml:"max_conns_per_host" json:"max_conns_per_host"`
		IdleConnTimeout       string `yaml:"idle_conn_timeout" json:"idle_conn_timeout"`
		ResponseHeaderTimeout string `yaml:"response_header_timeout" json:"response_header_timeout"`
		ExpectContinueTimeout string `yaml:"expect_continue_timeout" json:"expect_continue_timeout"`
		MaxResponseHeaderSize int64  `yaml:"max_response_header_size" json:"max_response_header_size"`
		WriteBufferSize       int64  `yaml:"write_buffer_size" json:"write_buffer_size"`
		ReadBufferSize        int64  `yaml:"read_buffer_size" json:"read_buffer_size"`
		ForceAttemptHTTP2     bool   `yaml:"force_attempt_http_2" json:"force_attempt_http_2"`
	} `yaml:"round_tripper" json:"round_tripper"`
	Backoff struct {
		// Algorithm  string // linear backoff
		Factor     float64 `yaml:"factor" json:"factor"`
		RetryCount int     `yaml:"retry_count" json:"retry_count"`
		TimeLimit  string  `yaml:"time_limit" json:"time_limit"`
	} `yaml:"backoff" json:"backoff"`
}

func (t *Transport) Bind() *Transport {
	t.RoundTripper.TLSHandshakeTimeout = GetActualValue(t.RoundTripper.TLSHandshakeTimeout)
	t.RoundTripper.IdleConnTimeout = GetActualValue(t.RoundTripper.IdleConnTimeout)
	t.RoundTripper.ResponseHeaderTimeout = GetActualValue(t.RoundTripper.ResponseHeaderTimeout)
	t.RoundTripper.ExpectContinueTimeout = GetActualValue(t.RoundTripper.ExpectContinueTimeout)
	t.Backoff.TimeLimit = GetActualValue(t.Backoff.TimeLimit)
	return t
}
