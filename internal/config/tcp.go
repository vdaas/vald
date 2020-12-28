//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

import (
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/tls"
)

// TCP represent the TCP configuration for server.
type TCP struct {
	DNS    *DNS    `yaml:"dns" json:"dns"`
	Dialer *Dialer `yaml:"dialer" json:"dialer"`
	TLS    *TLS    `yaml:"tls" json:"tls"`
}

type Dialer struct {
	Timeout          string `yaml:"timeout" json:"timeout"`
	KeepAlive        string `yaml:"keep_alive" json:"keep_alive"`
	DualStackEnabled bool   `yaml:"dual_stack_enabled" json:"dual_stack_enabled"`
}

type DNS struct {
	CacheEnabled    bool   `yaml:"cache_enabled" json:"cache_enabled"`
	RefreshDuration string `yaml:"refresh_duration" json:"refresh_duration"`
	CacheExpiration string `yaml:"cache_expiration" json:"cache_expiration"`
}

func (d *DNS) Bind() *DNS {
	d.RefreshDuration = GetActualValue(d.RefreshDuration)
	d.CacheExpiration = GetActualValue(d.CacheExpiration)
	return d
}

func (d *Dialer) Bind() *Dialer {
	d.Timeout = GetActualValue(d.Timeout)
	d.KeepAlive = GetActualValue(d.KeepAlive)
	return d
}

func (t *TCP) Bind() *TCP {
	if t.TLS != nil {
		t.TLS = t.TLS.Bind()
	}
	if t.DNS != nil {
		t.DNS = t.DNS.Bind()
	}
	if t.Dialer != nil {
		t.Dialer = t.Dialer.Bind()
	}
	return t
}

func (t *TCP) Opts() []tcp.DialerOption {
	opts := make([]tcp.DialerOption, 0, 7)
	if t.DNS != nil {
		opts = append(opts,
			tcp.WithDNSCacheExpiration(t.DNS.CacheExpiration),
			tcp.WithDNSRefreshDuration(t.DNS.RefreshDuration),
		)
		if t.DNS.CacheEnabled {
			opts = append(opts,
				tcp.WithEnableDNSCache(),
			)
		}
	}

	if t.Dialer != nil {
		opts = append(opts,
			tcp.WithDialerKeepAlive(t.Dialer.KeepAlive),
			tcp.WithDialerTimeout(t.Dialer.Timeout),
		)
		if t.Dialer.DualStackEnabled {
			opts = append(opts,
				tcp.WithEnableDialerDualStack(),
			)
		}
	}

	if t.TLS != nil && t.TLS.Enabled {
		cfg, err := tls.New(t.TLS.Opts()...)
		if err == nil {
			opts = append(opts,
				tcp.WithTLS(cfg),
			)
		}
	}

	return opts
}
