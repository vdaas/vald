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



// Package tcp provides tcp option
package tcp

import (
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/timeutil"
)

type DialerOption func(*dialer)

var (
	defaultDialerOptions = []DialerOption{
		WithDialerKeepAlive("30s"),
		WithDialerTimeout("30s"),
		WithEnableDialerDualStack(),
	}
)

func WithCache(c gache.Gache) DialerOption {
	return func(d *dialer) {
		d.cache = c
	}
}

func WithDNSRefreshDuration(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Minute * 30
		}
		d.dnsRefreshDuration = pd
	}
}

func WithDNSCacheExpiration(dur string) DialerOption {
	return func(d *dialer) {
		if dur == "" {
			return
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Hour
		}
		d.dnsCacheExpiration = pd
		if d.dnsCacheExpiration > 0 {
			WithEnableDNSCache()(d)
		}
	}
}

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

func WithDialerKeepAlive(dur string) DialerOption {
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

func WithEnableDNSCache() DialerOption {
	return func(d *dialer) {
		d.dnsCache = true
	}
}
func WithDisableDNSCache() DialerOption {
	return func(d *dialer) {
		d.dnsCache = false
	}
}

func WithEnableDialerDualStack() DialerOption {
	return func(d *dialer) {
		d.dialerDualStack = true
	}
}
func WithDisableDialerDualStack() DialerOption {
	return func(d *dialer) {
		d.dialerDualStack = false
	}
}
