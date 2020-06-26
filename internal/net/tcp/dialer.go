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
	"unsafe"

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
	ips    []string
	curIdx uint32
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

	if dc, err := d.lookup(dctx, host); err == nil {
		ipLen := len(dc.ips)

		switch ipLen {
		case 0:
		case 1:
			if conn, err := d.dial(dctx, network, dc.ips[0]+port); err == nil {
				return conn, nil
			}
			d.cache.Delete(host)
		default:
			// get and add the idx, and make sure another thread will get the different idx
			var idx int
			for {
				curIdx := atomic.LoadUint32(&dc.curIdx)
				idx = *(*int)(unsafe.Pointer(&curIdx))
				next := uint32((idx + 1) % ipLen)
				if atomic.CompareAndSwapUint32(&dc.curIdx, curIdx, next) {
					break
				}
			}

			// try the first dial and return if success
			if conn, err := d.dial(dctx, network, dc.ips[idx]+port); err == nil {
				return conn, nil
			}

			// if failed then try next and update the idx (not thread safe)
			for i := 1; i < ipLen; i++ {
				idx = (idx + i) % ipLen
				if conn, err := d.dial(dctx, network, dc.ips[idx]+port); err == nil {
					atomic.StoreUint32(&dc.curIdx, uint32(idx))
					return conn, nil
				}
			}

			// if all failed then remove the cache
			d.cache.Delete(host)
		}
	}

	return d.dial(dctx, network, addr)
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

	ips := make([]string, len(r))
	for i, ip := range r {
		ips[i] = ip.String()
	}

	dc := &dialerCache{
		ips: ips,
	}
	d.cache.Set(addr, dc)
	return dc, nil
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
