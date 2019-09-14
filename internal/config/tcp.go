//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//


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
