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
	"fmt"
	"math"
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

type dialerCache struct {
	ips []string
	cnt uint32
}

func (d *dialerCache) GetIP() string {
	if atomic.LoadUint32(&d.cnt) > math.MaxUint32/2 {
		atomic.StoreUint32(&d.cnt, 0)
	}
	return d.ips[atomic.AddUint32(&d.cnt, 1)%d.Len()]
}

func (d *dialerCache) Len() uint32 {
	return uint32(len(d.ips))
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

func (d *dialer) lookup(ctx context.Context, addr string) (*dialerCache, error) {
	cache, ok := d.cache.Get(addr)
	if ok {
		return cache.(*dialerCache), nil
	}

	r, err := d.der.Resolver.LookupIPAddr(ctx, addr)
	if err != nil {
		return nil, err
	}
	dc := &dialerCache{
		ips: make([]string, 0, len(r)),
	}
	for _, ip := range r {
		dc.ips = append(dc.ips, ip.String())
	}

	d.cache.Set(addr, dc)
	return dc, nil
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
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	if d.dnsCache {
		if dc, err := d.lookup(dctx, host); err == nil {
			for i := uint32(0); i < dc.Len(); i++ {
				if conn, err := d.dial(dctx, network, fmt.Sprintf("%s:%d", dc.GetIP(), port)); err == nil {
					return conn, nil
				}
			}
			d.cache.Delete(host)
		}
	}
	return d.dial(dctx, network, addr)
}

func (d *dialer) dial(ctx context.Context, network string, addr string) (net.Conn, error) {
	conn, err := d.der.DialContext(ctx, network, addr)
	if err != nil {
		defer func(conn net.Conn) {
			if conn != nil {
				conn.Close()
			}
		}(conn)
		return nil, err
	}

	if d.tlsConfig != nil {
		return tls.Client(conn, d.tlsConfig), nil
	}
	return conn, nil
}
