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
