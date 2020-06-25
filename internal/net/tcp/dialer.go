//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package tcp

import (
	"context"
	"crypto/tls"
	"strings"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
)

// Dialer is an interface to get the dialer instance to connect to an address.
type Dialer interface {
	GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error)
	StartDialerCache(ctx context.Context)
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type dialer struct {
	cache                 cache.Cache
	dnsCache              bool
	tlsConfig             *tls.Config
	dnsRefreshDurationStr string
	dnsCacheExpirationStr string
	dnsRefreshDuration    time.Duration
	dnsCacheExpiration    time.Duration
	dialerTimeout         time.Duration
	dialerKeepAlive       time.Duration
	dialerDualStack       bool
	der                   *net.Dialer
	dialer                func(ctx context.Context, network, addr string) (net.Conn, error)
	reqCnt                uint64
}

func NewDialer(opts ...DialerOption) (der Dialer, err error) {
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

	if !d.dnsCache {
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
		d.der.Resolver = &net.Resolver{
			PreferGo: false,
			Dial:     d.dialer,
		}
		return d, nil
	}

	if d.dnsRefreshDuration > d.dnsCacheExpiration {
		return nil, errors.ErrInvalidDNSConfig(d.dnsRefreshDuration, d.dnsCacheExpiration)
	}

	if d.cache == nil {
		d.cache, err = cache.New(
			cache.WithExpireDuration(d.dnsCacheExpirationStr),
			cache.WithExpireCheckDuration(d.dnsRefreshDurationStr),
			cache.WithExpiredHook(func(ctx context.Context, addr string) {
				if err := safety.RecoverFunc(func() (err error) {
					_, err = d.lookup(ctx, addr)
					return err
				}); err != nil {
					log.Error(err)
				}
			}),
		)
		if err != nil {
			return nil, err
		}
	}

	d.dialer = d.cachedDialer

	d.der.Resolver = &net.Resolver{
		PreferGo: false,
		Dial:     d.dialer,
	}

	return d, nil
}

func (d *dialer) GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.dialer
}

func (d *dialer) lookup(ctx context.Context, addr string) (ips []string, err error) {
	cache, ok := d.cache.Get(addr)
	if ok {
		return cache.([]string), nil
	}

	r, err := d.der.Resolver.LookupIPAddr(ctx, addr)
	if err != nil {
		return nil, err
	}

	ips = make([]string, len(r))
	for i, ip := range r {
		ips[i] = ip.String()
	}

	d.cache.Set(addr, ips)
	return ips, nil
}

func (d *dialer) StartDialerCache(ctx context.Context) {
	if d.dnsCache && d.cache != nil {
		d.cache.Start(ctx)
	}
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.GetDialer()(ctx, network, address)
}

func (d *dialer) cachedDialer(dctx context.Context, network, addr string) (conn net.Conn, err error) {
	sep := strings.LastIndex(addr, ":")
	if sep < 0 {
		sep = len(addr)
	}
	host := addr[:sep]
	port := addr[sep:] // includes ":" char

	if ips, err := d.lookup(dctx, host); err == nil {
		cnt := int(d.getAndAddReqCnt()) % len(ips)

		for i := cnt; i < len(ips); i++ {
			if conn, err := d.dial(dctx, network, ips[i]+port); err == nil {
				return conn, err
			}
		}
		if cnt > 1 {
			for i := 0; i < cnt; i++ {
				if conn, err := d.dial(dctx, network, ips[i]+port); err == nil {
					return conn, err
				}
			}
		}
		d.cache.Delete(host)
	}

	// TODO should use lookup
	return d.dial(dctx, network, addr)
}

func (d *dialer) dial(ctx context.Context, network string, addr string) (net.Conn, error) {
	conn, err := d.der.DialContext(ctx, network, addr)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return nil, err
	}

	if d.tlsConfig != nil {
		return tls.Client(conn, d.tlsConfig), nil
	}
	return conn, nil
}

func (d *dialer) getAndAddReqCnt() uint64 {
	for {
		i := atomic.LoadUint64(&d.reqCnt)
		if atomic.CompareAndSwapUint64(&d.reqCnt, i, i+1) {
			return i
		}
	}
}
