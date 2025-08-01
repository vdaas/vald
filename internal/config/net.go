//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

import (
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/tls"
)

// Net represents the network configuration tcp, udp, unix domain socket.
type Net struct {
	Network      string        `json:"network,omitempty"       yaml:"network"`
	DNS          *DNS          `json:"dns,omitempty"           yaml:"dns"`
	Dialer       *Dialer       `json:"dialer,omitempty"        yaml:"dialer"`
	SocketOption *SocketOption `json:"socket_option,omitempty" yaml:"socket_option"`
	TLS          *TLS          `json:"tls,omitempty"           yaml:"tls"`
}

// Dialer represents the configuration for dial.
type Dialer struct {
	Timeout          string `json:"timeout,omitempty"            yaml:"timeout"`
	Keepalive        string `json:"keepalive,omitempty"          yaml:"keepalive"`
	FallbackDelay    string `json:"fallback_delay,omitempty"     yaml:"fallback_delay"`
	DualStackEnabled bool   `json:"dual_stack_enabled,omitempty" yaml:"dual_stack_enabled"`
}

// DNS represents the configuration for resolving DNS.
type DNS struct {
	CacheEnabled    bool   `json:"cache_enabled,omitempty"    yaml:"cache_enabled"`
	RefreshDuration string `json:"refresh_duration,omitempty" yaml:"refresh_duration"`
	CacheExpiration string `json:"cache_expiration,omitempty" yaml:"cache_expiration"`
}

// SocketOption represents the socket configurations.
type SocketOption struct {
	ReusePort                bool `json:"reuse_port,omitempty"                  yaml:"reuse_port"`
	ReuseAddr                bool `json:"reuse_addr,omitempty"                  yaml:"reuse_addr"`
	TCPFastOpen              bool `json:"tcp_fast_open,omitempty"               yaml:"tcp_fast_open"`
	TCPNoDelay               bool `json:"tcp_no_delay,omitempty"                yaml:"tcp_no_delay"`
	TCPCork                  bool `json:"tcp_cork,omitempty"                    yaml:"tcp_cork"`
	TCPQuickAck              bool `json:"tcp_quick_ack,omitempty"               yaml:"tcp_quick_ack"`
	TCPDeferAccept           bool `json:"tcp_defer_accept,omitempty"            yaml:"tcp_defer_accept"`
	IPTransparent            bool `json:"ip_transparent,omitempty"              yaml:"ip_transparent"`
	IPRecoverDestinationAddr bool `json:"ip_recover_destination_addr,omitempty" yaml:"ip_recover_destination_addr"`
}

// Bind binds the actual data from the DNS fields.
func (d *DNS) Bind() *DNS {
	d.RefreshDuration = GetActualValue(d.RefreshDuration)
	d.CacheExpiration = GetActualValue(d.CacheExpiration)
	return d
}

// Bind binds the actual data from the Dialer fields.
func (d *Dialer) Bind() *Dialer {
	d.Timeout = GetActualValue(d.Timeout)
	d.Keepalive = GetActualValue(d.Keepalive)
	d.FallbackDelay = GetActualValue(d.FallbackDelay)
	return d
}

// Bind binds the actual data from the SocketOption fields.
func (s *SocketOption) Bind() *SocketOption {
	return s
}

// ToSocketFlag returns the control.SocketFlag defined as uint along with the SocketOption's fields.
func (s *SocketOption) ToSocketFlag() control.SocketFlag {
	var flg control.SocketFlag
	if s == nil {
		return flg
	}
	if s.ReuseAddr {
		flg |= control.ReuseAddr
	}
	if s.ReusePort {
		flg |= control.ReusePort
	}
	if s.TCPFastOpen {
		flg |= control.TCPFastOpen
	}
	if s.TCPCork {
		flg |= control.TCPCork
	}
	if s.TCPNoDelay {
		flg |= control.TCPNoDelay
	}
	if s.TCPDeferAccept {
		flg |= control.TCPDeferAccept
	}
	if s.TCPQuickAck {
		flg |= control.TCPQuickAck
	}
	if s.IPTransparent {
		flg |= control.IPTransparent
	}
	if s.IPRecoverDestinationAddr {
		flg |= control.IPRecoverDestinationAddr
	}
	return flg
}

// Bind binds the actual data from the Net fields.
func (t *Net) Bind() *Net {
	t.Network = GetActualValue(t.Network)
	if t.TLS != nil {
		t.TLS = t.TLS.Bind()
	}
	if t.DNS != nil {
		t.DNS = t.DNS.Bind()
	}
	if t.Dialer != nil {
		t.Dialer = t.Dialer.Bind()
	}
	if t.SocketOption != nil {
		t.SocketOption = t.SocketOption.Bind()
	}
	return t
}

// Opts creates the slice with the functional options for the net.Dialer options.
func (t *Net) Opts() ([]net.DialerOption, error) {
	opts := make([]net.DialerOption, 0, 7)
	if t.DNS != nil {
		opts = append(opts,
			net.WithDNSCacheExpiration(t.DNS.CacheExpiration),
			net.WithDNSRefreshDuration(t.DNS.RefreshDuration),
		)
		if t.DNS.CacheEnabled {
			opts = append(opts,
				net.WithEnableDNSCache(),
			)
		}
	}

	if t.Dialer != nil {
		opts = append(opts,
			net.WithDialerKeepalive(t.Dialer.Keepalive),
			net.WithDialerTimeout(t.Dialer.Timeout),
			net.WithDialerFallbackDelay(t.Dialer.FallbackDelay),
		)
		if t.Dialer.DualStackEnabled {
			opts = append(opts,
				net.WithEnableDialerDualStack(),
			)
		}
	}
	if t.SocketOption != nil {
		opts = append(opts, net.WithSocketFlag(t.SocketOption.ToSocketFlag()))
	}

	if t.TLS != nil && t.TLS.Enabled {
		cfg, err := tls.NewClientConfig(t.TLS.Opts()...)
		if err != nil {
			return nil, err
		}
		opts = append(opts,
			net.WithTLS(cfg),
		)
	}

	return opts, nil
}
