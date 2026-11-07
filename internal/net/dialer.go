// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"context"
	"net"
	"net/netip"
	"strconv"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/net/quic"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/tls"
)

// Dialer is an interface to get the dialer instance to connect to an address.
type Dialer interface {
	GetDialer() func(ctx context.Context, network, addr string) (Conn, error)
	StartDialerCache(ctx context.Context)
	DialContext(ctx context.Context, network, address string) (Conn, error)
}

type dialer struct {
	dnsCache              cacher.Cache[*dialerCache]
	ctrl                  control.SocketController
	dialer                func(ctx context.Context, network, addr string) (Conn, error)
	der                   *net.Dialer
	tlsConfig             *tls.Config
	addrs                 sync.Map[string, *addrInfo]
	dnsRefreshDurationStr string
	dnsCacheExpirationStr string
	dnsCacheExpiration    time.Duration
	dialerTimeout         time.Duration
	dialerKeepalive       time.Duration
	dialerFallbackDelay   time.Duration
	dnsRefreshDuration    time.Duration
	sockFlg               control.SocketFlag
	dnsCachedOnce         sync.Once
	dialerDualStack       bool
	enableDNSCache        bool
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

const (
	apiName = "vald/internal/net"

	// dnsRetryInterval is the pause between DNS lookup/dial retry rounds in cachedDialer.
	dnsRetryInterval = 10 * time.Millisecond
)

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
	d.ctrl = control.New(d.sockFlg, control.KeepAliveSeconds(d.dialerKeepalive))

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
				cache.WithExpireDuration[*dialerCache](d.dnsCacheExpirationStr),
				cache.WithExpireCheckDuration[*dialerCache](d.dnsRefreshDurationStr),
				cache.WithExpiredHook(d.cacheExpireHook),
			); err != nil {
				return nil, err
			}
		}
		d.dialer = d.cachedDialer
	}

	// The resolver dials DNS servers (udp, and tcp on truncation); those dials must
	// never be TLS-handshaked or go over QUIC, so it does not use d.dialer.
	d.der.Resolver = &Resolver{
		PreferGo: true,
		Dial:     d.resolverDial,
	}

	return d, nil
}

// GetDialer returns a function to return the connection.
func (d *dialer) GetDialer() func(ctx context.Context, network, addr string) (Conn, error) {
	return d.dialer
}

func (d *dialer) lookup(ctx context.Context, host string) (dc *dialerCache, err error) {
	if d.enableDNSCache {
		if dc, ok := d.dnsCache.Get(host); ok {
			if dc != nil && len(dc.ips) > 0 {
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
		if d.enableDNSCache {
			d.dnsCache.Delete(host)
		}
		return nil, err
	}
	if len(ips) == 0 {
		if d.enableDNSCache {
			d.dnsCache.Delete(host)
		}
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
	if d.der != nil && d.der.Resolver != nil {
		rsv = d.der.Resolver
	}
	r, err := lookupNetIP(ctx, rsv, host)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrLookupIPAddrNotFound(host).Error())
	}
	if len(r) == 0 {
		return nil, errors.ErrLookupIPAddrNotFound(host)
	}
	ips = make([]string, 0, len(r))
	for _, ip := range r {
		if ip.Is4In6() {
			ip = ip.Unmap()
		}
		ips = append(ips, ip.StringExpanded())
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
		host, nport, _, isV4, isV6, err = Parse(ctx, addr)
		if err != nil {
			if d.enableDNSCache {
				idx := strings.LastIndex(addr, ":")
				if idx >= 0 {
					host = addr[:idx]
				} else {
					host = addr
				}
				d.dnsCache.Delete(host)
			}
			return nil, err
		}
		// Parse substitutes defaultPort for a missing port, so nport is never 0 here.
		port = strconv.FormatUint(uint64(nport), 10)
		if _, perr := netip.ParseAddr(host); perr == nil {
			// An IP literal dials itself and re-parsing it needs no DNS lookup, so it
			// is neither cached here nor sent through the DNS cache. gRPC hands the
			// dialer every resolved IP:port, and caching them kept one entry per
			// backend IP ever dialed: Disconnect evicts only the configured address.
			isIP = true
		} else {
			// Keep the upstream semantics: Parse also reports isIPv4/isIPv6 for resolved
			// host names, so later dials skip the DNS-cache path while the first one
			// (local isIP still false) goes through it. Making host names use the cache
			// path on every dial changes TLS server names and retry timing.
			d.addrs.Store(addr, &addrInfo{
				host: host,
				port: port,
				addr: addr,
				isIP: isV4 || isV6,
			})
		}
	} else {
		host = ai.host
		port = ai.port
		isIP = ai.isIP
	}

	if d.enableDNSCache && !isIP {
		to := time.NewTimer(d.dialerTimeout)
		defer to.Stop()
		for {
			select {
			case <-to.C:
				d.dnsCache.Delete(host)
				conn, err = d.dial(ctx, network, addr)
				if err != nil {
					d.dnsCache.Delete(host)
					return nil, err
				}
				return conn, nil
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
							// Verify TLS against the host name, not the resolved IP that
							// target carries, as the plain dial of addr does.
							tconn, err := d.dialConn(ctx, network, target, host, true)
							if err == nil && tconn != nil {
								return tconn, nil
							}
							log.Warnf("failed to dial connection to %s\terror: %v", target, err)
							if tconn != nil {
								if cerr := tconn.Close(); cerr != nil {
									log.Warnf("failed to close connection to %s\terror: %v", target, cerr)
								}
							}
						}
						d.dnsCache.Delete(host)
					}
				}
				// Back off before retrying so that a fast-failing lookup/dial does not spin
				// the CPU and the resolver until dialerTimeout, and so that a cancelled ctx
				// ends the loop immediately.
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(dnsRetryInterval):
				}
			}
		}
	}
	conn, err = d.dial(ctx, network, addr)
	if err != nil && d.enableDNSCache {
		d.dnsCache.Delete(host)
		return nil, err
	}
	return conn, err
}

func isQUICDial(network, addr string) bool {
	if !IsUDP(network) {
		return false
	}
	host, port, err := SplitHostPort(addr)
	if err != nil || host == "" || port == 0 {
		return false
	}
	return port != 53
}

func (d *dialer) dial(ctx context.Context, network, addr string) (conn Conn, err error) {
	return d.dialConn(ctx, network, addr, "", true)
}

// resolverDial is the Resolver's dial func: a plain connection, never TLS or QUIC.
func (d *dialer) resolverDial(ctx context.Context, network, addr string) (conn Conn, err error) {
	return d.dialConn(ctx, network, addr, "", false)
}

// dialConn dials addr. serverName, when set, is the TLS server name to verify
// instead of the host of addr.
func (d *dialer) dialConn(
	ctx context.Context, network, addr, serverName string, secure bool,
) (conn Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.dial")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if NetworkTypeFromString(network) == Unknown {
		network = TCP.String()
	}
	if addr == "" {
		return nil, errors.ErrInvalidAddress(network, addr)
	}
	log.Debugf("%s connection dialing to addr %s", network, addr)
	// Capture this before the closure below rewrites network from udp to tcp, so
	// that udp dials (e.g. the resolver's DNS dials) are never TLS-handshaked.
	isUDP := IsUDP(network)
	err = safety.RecoverWithoutPanicFunc(func() (err error) {
		if secure && isQUICDial(network, addr) {
			conn, err = quic.DialContext(ctx, addr, d.tlsConfig)
		} else {
			if IsUDP(network) {
				network = TCP.String()
			}
			conn, err = d.der.DialContext(ctx, network, addr)
		}
		return err
	})()
	if err != nil {
		if conn != nil {
			err = errors.Join(conn.Close(), err)
		}
		return nil, err
	}

	if secure && !isUDP && d.tlsConfig != nil {
		// Do not return the *tls.Conn directly: a nil *tls.Conn converted to Conn is a
		// non-nil interface, and callers that Close a non-nil conn on error would panic.
		tconn, err := d.tlsHandshake(ctx, conn, network, addr, serverName)
		if err != nil {
			return nil, err
		}
		return tconn, nil
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

func (d *dialer) tlsHandshake(
	ctx context.Context, conn Conn, network, addr, serverName string,
) (tconn *tls.Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Dialer.tlsHandshake")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	// d.tlsConfig is shared by every dial (possibly to different hosts), so it must
	// never be mutated here: work on a per-dial clone to set the ServerName.
	cfg := d.tlsConfig.Clone()
	if cfg.ServerName == "" && serverName != "" {
		cfg.ServerName = serverName
	}
	if cfg.ServerName == "" {
		var host string
		host, _, err = SplitHostPort(addr)
		if err == nil && len(host) != 0 {
			cfg.ServerName = host
		}
	}
	if conn != nil {
		tconn = tls.Client(conn, cfg)
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
		// The in-place handshake on the passed-in conn failed; close it so its
		// socket is not leaked before we re-dial a fresh connection below.
		if tconn != nil {
			if cerr := tconn.Close(); cerr != nil {
				log.Debugf("failed to close failed-handshake tls conn for %s: %v", addr, cerr)
			}
			tconn = nil
		}
		// Only bound the re-dial when a positive timeout is configured: a zero
		// d.der.Timeout means "no timeout", and WithTimeout(ctx, 0) would hand
		// DialContext an already-expired context (mirrors the primary attempt).
		tctx := ctx
		if d.der.Timeout > 0 {
			var tcancel context.CancelFunc
			tctx, tcancel = context.WithTimeout(ctx, d.der.Timeout)
			defer tcancel()
		}
		err = safety.RecoverWithoutPanicFunc(func() error {
			tder := &tls.Dialer{
				NetDialer: d.der,
				Config:    cfg,
			}
			conn, err = tder.DialContext(tctx, network, addr)
			return err
		})()
		if err != nil || conn == nil {
			ttctx := ctx
			if d.der.Timeout > 0 {
				var ttcancel context.CancelFunc
				ttctx, ttcancel = context.WithTimeout(ctx, d.der.Timeout)
				defer ttcancel()
			}
			err = safety.RecoverWithoutPanicFunc(func() error {
				tder := &tls.Dialer{
					Config: cfg,
				}
				conn, err = tder.DialContext(ttctx, network, addr)
				return err
			})()
		}
		if err != nil || conn == nil {
			if conn != nil {
				if err != nil {
					return nil, errors.Join(conn.Close(), err)
				}
				return nil, conn.Close()
			}
			return nil, err
		}
		// Assign to the named return tconn instead of shadowing it with :=, so
		// the successfully re-dialed connection (not the original failed one
		// from the top of the function) is what gets returned below.
		var ok bool
		tconn, ok = conn.(*tls.Conn)
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

// EvictAddr drops the cached parse result of addr and the DNS cache entry of its
// host, so the state of removed addresses does not pile up. The gRPC client
// calls it when it disconnects addr. A later dial to addr re-creates both.
func (d *dialer) EvictAddr(addr string) {
	ai, ok := d.addrs.LoadAndDelete(addr)
	if ok && ai != nil && d.dnsCache != nil {
		d.dnsCache.Delete(ai.host)
	}
}

func (d *dialer) cacheExpireHook(ctx context.Context, addr string, dc *dialerCache) {
	// An IP literal always resolves to itself, so refreshing it would keep its
	// entry forever. cachedDialer never reads one (IP literals skip the cache),
	// so let a stored one expire instead.
	if _, err := netip.ParseAddr(addr); err == nil {
		return
	}
	if err := safety.RecoverFunc(func() (err error) {
		dc, err = d.lookup(ctx, addr)
		return err
	})(); err != nil {
		if dc != nil {
			log.Errorf("dns cache expiration hook process returned error: %v\tfor addr:\t%s\tips: %v\tlen: %d", err, addr, dc.ips, dc.Len())
		} else {
			log.Errorf("dns cache expiration hook process returned error: %v\tfor addr:\t%s", err, addr)
		}
	}
}
