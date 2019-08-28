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
type TCP struct {
	DNS struct {
		CacheEnabled    bool   `yaml:"cache_enabled" json:"cache_enabled"`
		RefreshDuration string `yaml:"refresh_duration" json:"refresh_duration"`
		CacheExpiration string `yaml:"cache_expiration" json:"cache_expiration"`
	} `yaml:"dns" json:"dns"`
	Dialer struct {
		Timeout          string `yaml:"timeout" json:"timeout"`
		KeepAlive        string `yaml:"keep_alive" json:"keep_alive"`
		DualStackEnabled bool   `yaml:"dual_stack_enabled" json:"dual_stack_enabled"`
	} `yaml:"dialer" json:"dialer"`
}

func (t *TCP) Bind() *TCP {
	t.DNS.RefreshDuration = GetActualValue(t.DNS.RefreshDuration)
	t.DNS.CacheExpiration = GetActualValue(t.DNS.CacheExpiration)
	t.Dialer.Timeout = GetActualValue(t.Dialer.Timeout)
	t.Dialer.KeepAlive = GetActualValue(t.Dialer.KeepAlive)
	return t
}
