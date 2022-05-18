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

	"github.com/cloudwego/netpoll"
	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/net/quic"
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
	cache                 cache.Cache
	dnsCache              bool
	dnsCachedOnce         sync.Once
	tlsConfig             *tls.Config
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
	npDialer              netpoll.Dialer
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
	netpoll.SetLoadBalance(netpoll.RoundRobin)
	d.npDialer = netpoll.NewDialer()

	d.dialer = d.dial

	if d.dnsCache {
		if d.dnsRefreshDuration > d.dnsCacheExpiration {
			return nil, errors.ErrInvalidDNSConfig(d.dnsRefreshDuration, d.dnsCacheExpiration)
		}
		if d.cache == nil {
			if d.cache, err = cache.New(
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
		PreferGo: false,
		Dial:     d.dialer,
	}

	return d, nil
}

// GetDialer returns a function to return the connection.
func (d *dialer) GetDialer() func(ctx context.Context, network, addr string) (Conn, error) {
	return d.dialer
}

func (d *dialer) lookup(ctx context.Context, host string) (*dialerCache, error) {
	cache, ok := d.cache.Get(host)
	if ok {
		return cache.(*dialerCache), nil
	}
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.lookup")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	r, err := d.der.Resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}

	if len(r) == 0 {
		return nil, errors.ErrLookupIPAddrNotFound(host)
	}

	dc := &dialerCache{
		ips: make([]string, 0, len(r)),
	}

	for _, ip := range r {
		dc.ips = append(dc.ips, ip.String())
	}

	if dc != nil && len(dc.ips) != 0 {
		log.Infof("lookup succeed %v", dc.ips)
		d.cache.Set(host, dc)
	}
	return dc, nil
}

// StartDialerCache starts the dialer cache to expire the cache automatically.
func (d *dialer) StartDialerCache(ctx context.Context) {
	if d.dnsCache && d.cache != nil {
		d.dnsCachedOnce.Do(func() {
			d.cache.Start(ctx)
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
			d.addrs.Delete(addr)
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

	if d.dnsCache && !isIP {
		if dc, err := d.lookup(ctx, host); err == nil {
			for i := uint32(0); i < dc.Len(); i++ {
				// in this line we use golang's standard net packages net.JoinHostPort cuz port is string type
				target := net.JoinHostPort(dc.IP(), port)
				conn, err := d.dial(ctx, network, target)
				if err == nil && conn != nil {
					return conn, nil
				}
				log.Warnf("failed to dial connection to %s\terror: %v", target, err)
			}
			d.cache.Delete(host)
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
	if IsUDP(network) {
		conn, err = quic.DialQuicContext(ctx, addr, d.tlsConfig)
	} else {
		conn, err = d.npDialer.DialConnection(network, addr, d.der.Timeout)
		if err != nil {
			conn, err = d.der.DialContext(ctx, network, addr)
		}
	}
	if err != nil {
		defer func(conn Conn) {
			if conn != nil {
				if err != nil {
					err = errors.Wrap(conn.Close(), err.Error())
					return
				}
				err = conn.Close()
			}
		}(conn)
		return nil, err
	}

	if !IsUDP(network) && d.tlsConfig != nil {
		return d.tlsHandshake(ctx, conn, addr)
	}
	if conn != nil {
		log.Infof("connected to addr %s succeed from %s://%s to %s://%s",
			addr,
			conn.LocalAddr().Network(), conn.LocalAddr().String(),
			conn.RemoteAddr().Network(), conn.RemoteAddr().String(),
		)
	}
	return conn, nil
}

func (d *dialer) tlsHandshake(ctx context.Context, conn Conn, addr string) (*tls.Conn, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.tlsHandshake")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var err error
	if d.tlsConfig.ServerName == "" {
		var host string
		host, _, err = SplitHostPort(addr)
		if err == nil {
			d.tlsConfig.ServerName = host
		}
	}
	tconn := tls.Client(conn, d.tlsConfig)
	var tctx context.Context
	if d.der.Timeout > 0 {
		var cancel context.CancelFunc
		tctx, cancel = context.WithTimeout(ctx, d.der.Timeout)
		defer cancel()
	} else {
		tctx = ctx
	}
	err = tconn.HandshakeContext(tctx)
	if err != nil {
		defer func(conn Conn) {
			if conn != nil {
				if err != nil {
					err = errors.Wrap(conn.Close(), err.Error())
					return
				}
				err = conn.Close()
			}
		}(conn)
		return nil, err
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
		log.Errorf("dns cacheExpireHook error occurred: %v\taddr:\t%s", err, addr)
	}
}
