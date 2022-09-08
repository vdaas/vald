//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package tcp provides tcp option
package net

import (
	"crypto/tls"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/timeutil"
)

// DialerOption represent the functional option for dialer.
type DialerOption func(*dialer)

var defaultDialerOptions = []DialerOption{
	WithDialerKeepalive("30s"),
	WithDialerTimeout("0s"),
	WithDialerFallbackDelay("300ms"),
	WithEnableDialerDualStack(),
	WithDisableDNSCache(),
}

// WithDNSCache returns the functional option to set the cache.
func WithDNSCache(c cache.Cache) DialerOption {
	return func(d *dialer) {
		d.dnsCache = c
		if d.dnsCache != nil {
			WithEnableDNSCache()(d)
		}
	}
}

// WithSocketFlag returns the functional option to set the socket options.
func WithSocketFlag(flg control.SocketFlag) DialerOption {
	return func(d *dialer) {
		d.sockFlg = flg
	}
}

// WithDNSRefreshDuration returns the functional option to set the DNSRefreshDuration.
func WithDNSRefreshDuration(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			WithDNSRefreshDuration("30m")(d)
			return
		}
		d.dnsRefreshDuration = pd
		d.dnsRefreshDurationStr = dur
	}
}

// WithDNSCacheExpiration returns the functional option to set the DNSCacheExpiration.
func WithDNSCacheExpiration(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			WithDNSCacheExpiration("1h")(d)
			return
		}
		d.dnsCacheExpiration = pd
		d.dnsCacheExpirationStr = dur
		if d.dnsCacheExpiration > 0 {
			WithEnableDNSCache()(d)
		}
	}
}

// WithDialerTimeout returns the functional option to set the DialerTimeout.
func WithDialerTimeout(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Second * 30
		}
		d.dialerTimeout = pd
	}
}

// WithDialerKeepalive returns the functional option to set the DialerKeepalive.
func WithDialerKeepalive(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Second * 30
		}
		d.dialerKeepalive = pd
	}
}

// WithDialerFallbackDelay returns the functional option to set the DialerFallbackDelay.
func WithDialerFallbackDelay(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Millisecond * 300
		}
		d.dialerFallbackDelay = pd
	}
}

// WithTLS returns the functional option to set the DialerTLS.
func WithTLS(cfg *tls.Config) DialerOption {
	return func(d *dialer) {
		d.tlsConfig = cfg
	}
}

// WithEnableDNSCache returns the functional option to enable DNSCache.
func WithEnableDNSCache() DialerOption {
	return func(d *dialer) {
		d.enableDNSCache = true
	}
}

// WithDisableDNSCache returns the functional option to disable DNSCache.
func WithDisableDNSCache() DialerOption {
	return func(d *dialer) {
		d.enableDNSCache = false
	}
}

// WithEnableDialerDualStack returns the functional option to enable DialerDualStack.
func WithEnableDialerDualStack() DialerOption {
	return func(d *dialer) {
		d.dialerDualStack = true
	}
}

// WithDisableDialerDualStack returns the functional option to disable DialerDualStack.
func WithDisableDialerDualStack() DialerOption {
	return func(d *dialer) {
		d.dialerDualStack = false
	}
}
