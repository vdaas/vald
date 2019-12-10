//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
	"context"
	"crypto/tls"
	"net"
	"strings"
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

// type Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

type Dialer interface {
	GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error)
	StartDialerCache(ctx context.Context)
}

type dialer struct {
	cache              gache.Gache
	dnsCache           bool
	tlsConfig          *tls.Config
	dnsRefreshDuration time.Duration
	dnsCacheExpiration time.Duration
	dialerTimeout      time.Duration
	dialerKeepAlive    time.Duration
	dialerDualStack    bool
	der                *net.Dialer
	dialer             func(ctx context.Context, network, addr string) (net.Conn, error)
}

func NewDialer(opts ...DialerOption) Dialer {
	d := new(dialer)
	for _, opt := range append(defaultDialerOptions, opts...) {
		opt(d)
	}

	d.der = &net.Dialer{
		Timeout:   d.dialerTimeout,
		KeepAlive: d.dialerKeepAlive,
		DualStack: d.dialerDualStack,
		Control:   Control,
	}

	d.der.Resolver = &net.Resolver{
		PreferGo: false,
		Dial:     d.der.DialContext,
	}

	if !d.dnsCache || d.cache == nil {
		if d.tlsConfig != nil {
			d.dialer = func(ctx context.Context, network,
				addr string) (conn net.Conn, err error) {
				conn, err = d.der.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return tls.Client(conn, d.tlsConfig), nil
			}
		} else {
			d.dialer = d.der.DialContext
		}
		return d
	}
	if d.cache == nil {
		d.cache = gache.New()
	}

	if d.dnsRefreshDuration > d.dnsCacheExpiration {
		d.dnsRefreshDuration, d.dnsCacheExpiration =
			d.dnsCacheExpiration, d.dnsRefreshDuration
	}

	d.dialer = d.cachedDialer

	return d
}

func (d *dialer) GetDialer() func(ctx context.Context,
	network, addr string) (net.Conn, error) {
	return d.dialer
}

func (d *dialer) lookup(ctx context.Context,
	addr string) (ips map[int]string, err error) {
	cache, ok := d.cache.Get(addr)
	if ok {
		return cache.(map[int]string), nil
	}

	r, err := d.der.Resolver.LookupIPAddr(ctx, addr)
	if err != nil {
		return nil, err
	}

	ips = make(map[int]string, len(r))
	for i, ip := range r {
		ips[i] = ip.String()
	}

	d.cache.SetWithExpire(addr, ips,
		d.dnsCacheExpiration)

	return ips, nil
}

func (d *dialer) StartDialerCache(ctx context.Context) {
	if d.dnsCache && d.cache != nil {
		d.cache.SetDefaultExpire(d.dnsCacheExpiration).
			SetExpiredHook(func(gctx context.Context, addr string) {
				if err := safety.RecoverFunc(func() (err error) {
					_, err = d.lookup(gctx, addr)
					return err
				}); err != nil {
					log.Error(err)
				}
			}).
			EnableExpiredHook().
			StartExpired(ctx, d.dnsRefreshDuration)
	}
}

func (d *dialer) cachedDialer(dctx context.Context, network, addr string) (
	conn net.Conn, err error) {
	sep := strings.LastIndex(addr, ":")

	if sep < 0 {
		sep = len(addr)
	}

	ips, err := d.lookup(dctx, addr[:sep])
	if err == nil {
		for _, ip := range ips {
			conn, err = d.der.DialContext(dctx, network, ip+addr[sep:])
			if err == nil {
				if d.tlsConfig != nil {
					return tls.Client(conn, d.tlsConfig), nil
				}
				return conn, nil
			}
			if conn != nil {
				conn.Close()
			}
		}
		d.cache.Delete(addr[:sep])
	}

	conn, err = d.der.DialContext(dctx, network, addr)
	if d.tlsConfig != nil {
		return tls.Client(conn, d.tlsConfig), nil
	}
	return
}
