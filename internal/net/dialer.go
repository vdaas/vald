//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package net provides net functionality for vald's network connection
package net

import (
	"context"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/tls"
)

// Dialer is an interface to get the dialer instance to connect to an address.
type Dialer interface {
	GetDialer() func(ctx context.Context, network, addr string) (Conn, error)
	StartDialerCache(ctx context.Context)
	DialContext(ctx context.Context, network, address string) (Conn, error)
}

type dialer struct {
	dnsCache              cacher.Cache
	enableDNSCache        bool
	dnsCachedOnce         sync.Once
	tlsConfig             *tls.Config
	tmu                   sync.RWMutex // lock mutex for tls handshake update
	dnsRefreshDurationStr string
	dnsCacheExpirationStr string
	dnsRefreshDuration    time.Duration
	dnsCacheExpiration    time.Duration
	dialerTimeout         time.Duration
	dialerKeepalive       time.Duration
	dialerFallbackDelay   time.Duration
	ctrl                  control.SocketController
	sockFlg               control.SocketFlag
	dialerDualStack       bool
	addrs                 sync.Map
	der                   *net.Dialer
	dialer                func(ctx context.Context, network, addr string) (Conn, error)
}

type addrInfo struct {
	addr string
	host string
	port string
	isIP bool
}

type dialerCache struct {
	ips []string
	cnt uint32
}

const apiName = "vald/internal/net"

// IP returns the next cached IP address in round robin order.
// It starts getting the index 1 cache instead of index 0 cache.
func (d *dialerCache) IP() string {
	if d.Len() == 1 {
		return d.ips[0]
	}

	return d.ips[atomic.AddUint32(&d.cnt, 1)%d.Len()]
}

// Len returns the length of cached IP addresses.
func (d *dialerCache) Len() uint32 {
	return uint32(len(d.ips))
}

// NewDialer initialize and return the dialer instance.
func NewDialer(opts ...DialerOption) (der Dialer, err error) {
	d := new(dialer)
	for _, opt := range append(defaultDialerOptions, opts...) {
		opt(d)
	}
	d.ctrl = control.New(d.sockFlg, int(d.dialerKeepalive))

	d.der = &net.Dialer{
		Timeout:       d.dialerTimeout,
		KeepAlive:     d.dialerKeepalive,
		DualStack:     d.dialerDualStack,
		FallbackDelay: d.dialerFallbackDelay,
		Control: func(network, addr string, c syscall.RawConn) (err error) {
			if d.ctrl != nil {
				return d.ctrl.GetControl()(network, addr, c)
			}
			return nil
		},
	}

	d.dialer = d.dial
	if d.enableDNSCache {
		if d.dnsRefreshDuration > d.dnsCacheExpiration {
			return nil, errors.ErrInvalidDNSConfig(d.dnsRefreshDuration, d.dnsCacheExpiration)
		}
		if d.dnsCache == nil {
			if d.dnsCache, err = cache.New(
				cache.WithExpireDuration(d.dnsCacheExpirationStr),
				cache.WithExpireCheckDuration(d.dnsRefreshDurationStr),
				cache.WithExpiredHook(d.cacheExpireHook),
			); err != nil {
				return nil, err
			}
		}
		d.dialer = d.cachedDialer
	}

	d.der.Resolver = &Resolver{
		PreferGo: true,
		Dial:     d.dialer,
	}

	return d, nil
}

// GetDialer returns a function to return the connection.
func (d *dialer) GetDialer() func(ctx context.Context, network, addr string) (Conn, error) {
	return d.dialer
}

func (d *dialer) lookup(ctx context.Context, host string) (dc *dialerCache, err error) {
	if d.enableDNSCache {
		dnsCache, ok := d.dnsCache.Get(host)
		if ok && dnsCache != nil {
			dc, ok = dnsCache.(*dialerCache)
			if ok && dc != nil && len(dc.ips) > 0 {
				return dc, nil
			}
		}
	}
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.lookup")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ips, err := d.lookupIPAddrs(ctx, host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.ErrLookupIPAddrNotFound(host)
	}

	dc = &dialerCache{
		ips: ips,
	}
	log.Debugf("lookup succeed for %s, ips: %v", host, dc.ips)
	if d.enableDNSCache {
		d.dnsCache.Set(host, dc)
	}
	return dc, nil
}

func (d *dialer) lookupIPAddrs(ctx context.Context, host string) (ips []string, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.lookupIPAddrs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var rsv *net.Resolver
	if d.der == nil || d.der.Resolver == nil {
		rsv = DefaultResolver
	} else {
		rsv = d.der.Resolver
	}

	r, err := rsv.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, errors.ErrLookupIPAddrNotFound(host)
	}

	ips = make([]string, 0, len(r))

	for _, ip := range r {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

// StartDialerCache starts the dialer cache to expire the cache automatically.
func (d *dialer) StartDialerCache(ctx context.Context) {
	if d.enableDNSCache && d.dnsCache != nil {
		d.dnsCachedOnce.Do(func() {
			d.dnsCache.Start(ctx)
		})
	}
}

// DialContext returns the connection or error base on the input.
// If the DNS cache is enabled, it will lookup the DNS cache in round robin order and return a connection of it.
// Also if TLS is enabled, it will create a TLS connection for it.
func (d *dialer) DialContext(ctx context.Context, network, address string) (Conn, error) {
	return d.GetDialer()(ctx, network, address)
}

func (d *dialer) cachedDialer(ctx context.Context, network, addr string) (conn Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.cachedDialer")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var (
		host string
		port string
		isIP bool
	)
	ai, ok := d.addrs.Load(addr)
	if !ok {
		var nport uint16
		var isV4, isV6 bool
		host, nport, _, isV4, isV6, err = Parse(addr)
		if err != nil {
			return nil, err
		}
		port = strconv.FormatUint(uint64(nport), 10)
		d.addrs.Store(addr, &addrInfo{
			host: host,
			port: port,
			addr: addr,
			isIP: isV4 || isV6,
		})
	} else {
		info, ok := ai.(*addrInfo)
		if ok {
			host = info.host
			port = info.port
			isIP = info.isIP
		}
	}

	if d.enableDNSCache && !isIP {
		to := time.NewTimer(d.dialerTimeout)
		defer to.Stop()
		for {
			select {
			case <-to.C:
				d.dnsCache.Delete(host)
				return d.dial(ctx, network, addr)
			default:
				if dc, err := d.lookup(ctx, host); err == nil {
					for i := uint32(0); i < dc.Len(); i++ {
						select {
						case <-to.C:
							d.dnsCache.Delete(host)
							return d.dial(ctx, network, addr)
						default:
							// in this line we use golang's standard net packages net.JoinHostPort cuz port is string type
							target := net.JoinHostPort(dc.IP(), port)
							conn, err := d.dial(ctx, network, target)
							if err == nil && conn != nil {
								return conn, nil
							}
							log.Warnf("failed to dial connection to %s\terror: %v", target, err)
							if conn != nil {
								conn.Close()
							}
						}
						d.dnsCache.Delete(host)
					}
				}
			}
		}
	}
	return d.dial(ctx, network, addr)
}

func (d *dialer) dial(ctx context.Context, network, addr string) (conn Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.dial")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	log.Debugf("%s connection dialing to addr %s", network, addr)
	err = safety.RecoverWithoutPanicFunc(func() error {
		conn, err = d.der.DialContext(ctx, network, addr)
		return err
	})()
	if err != nil {
		defer func(conn Conn) {
			if conn != nil {
				if err != nil {
					err = errors.Join(conn.Close(), err)
					return
				}
				err = conn.Close()
			}
		}(conn)
		return nil, err
	}

	d.tmu.RLock()
	if d.tlsConfig != nil {
		d.tmu.RUnlock()
		return d.tlsHandshake(ctx, conn, network, addr)
	}
	d.tmu.RUnlock()
	if conn != nil {
		log.Infof("connected to addr %s succeed from %s://%s to %s://%s",
			addr,
			conn.LocalAddr().Network(), conn.LocalAddr().String(),
			conn.RemoteAddr().Network(), conn.RemoteAddr().String(),
		)
	}
	return conn, nil
}

func (d *dialer) tlsHandshake(ctx context.Context, conn Conn, network, addr string) (tconn *tls.Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.tlsHandshake")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	d.tmu.RLock()
	if d.tlsConfig.ServerName == "" {
		d.tmu.RUnlock()
		var host string
		host, _, err = SplitHostPort(addr)
		if err == nil && len(host) != 0 {
			d.tmu.Lock()
			d.tlsConfig.ServerName = host
			d.tmu.Unlock()
		}
	} else {
		d.tmu.RUnlock()
	}
	if conn != nil {
		d.tmu.RLock()
		tconn = tls.Client(conn, d.tlsConfig)
		d.tmu.RUnlock()
	}
	var tctx context.Context
	if d.der.Timeout > 0 {
		var cancel context.CancelFunc
		tctx, cancel = context.WithTimeout(ctx, d.der.Timeout)
		defer cancel()
	} else {
		tctx = ctx
	}
	if tconn != nil {
		err = safety.RecoverWithoutPanicFunc(func() error {
			return tconn.HandshakeContext(tctx)
		})()
		if err == nil && !tconn.ConnectionState().HandshakeComplete {
			err = errors.ErrFailedToHandshakeTLSConnection(network, addr)
		}
	} else {
		err = errors.ErrFailedToHandshakeTLSConnection(network, addr)
	}
	if err != nil {
		tctx, tcancel := context.WithTimeout(ctx, d.der.Timeout)
		defer tcancel()
		err = safety.RecoverWithoutPanicFunc(func() error {
			d.tmu.RLock()
			tder := &tls.Dialer{
				NetDialer: d.der,
				Config:    d.tlsConfig,
			}
			d.tmu.RUnlock()
			conn, err = tder.DialContext(tctx, network, addr)
			return err
		})()
		if err != nil || conn == nil {
			ttctx, ttcancel := context.WithTimeout(ctx, d.der.Timeout)
			defer ttcancel()
			err = safety.RecoverWithoutPanicFunc(func() error {
				d.tmu.RLock()
				tder := &tls.Dialer{
					Config: d.tlsConfig,
				}
				d.tmu.RUnlock()
				conn, err = tder.DialContext(ttctx, network, addr)
				return err
			})()
		}
		if err != nil || conn == nil {
			defer func(conn Conn) {
				if conn != nil {
					if err != nil {
						err = errors.Join(conn.Close(), err)
						return
					}
					err = conn.Close()
				}
			}(conn)
			return nil, err
		}
		tconn, ok := conn.(*tls.Conn)
		if !ok || tconn == nil || !tconn.ConnectionState().HandshakeComplete {
			return nil, errors.ErrFailedToHandshakeTLSConnection(network, addr)
		}
	}
	if tconn != nil {
		log.Infof("tls handshake addr %s succeed from %s://%s to %s://%s,\tconnectionstate: [ Version:%d, ServerName: %s, HandshakeComplete: %v, DidResume: %v, NegotiatedProtocol: %s ]",
			addr,
			tconn.LocalAddr().Network(), tconn.LocalAddr().String(),
			tconn.RemoteAddr().Network(), tconn.RemoteAddr().String(),
			tconn.ConnectionState().Version,
			tconn.ConnectionState().ServerName,
			tconn.ConnectionState().HandshakeComplete,
			tconn.ConnectionState().DidResume,
			tconn.ConnectionState().NegotiatedProtocol,
		)
	}
	return tconn, nil
}

func (d *dialer) cacheExpireHook(ctx context.Context, addr string) {
	if err := safety.RecoverFunc(func() (err error) {
		_, err = d.lookup(ctx, addr)
		return
	})(); err != nil {
		log.Errorf("dns cache expiration hook process returned error: %v\tfor addr:\t%s", err, addr)
	}
}
