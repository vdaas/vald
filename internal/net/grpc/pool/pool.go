//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package pool provides gRPC connection pool client
package pool

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type (
	ClientConn = grpc.ClientConn
	DialOption = grpc.DialOption
)

type Conn interface {
	Connect(context.Context) (Conn, error)
	Disconnect(context.Context) error
	Do(ctx context.Context, f func(*ClientConn) error) error
	Get(context.Context) (conn *ClientConn, ok bool)
	IsHealthy(context.Context) bool
	IsIPConn() bool
	Len() uint64
	Size() uint64
	Reconnect(ctx context.Context, force bool) (Conn, error)
	String() string
}

type poolConn struct {
	conn *ClientConn
	addr string
}

type pool struct {
	pool          []atomic.Pointer[poolConn]
	startPort     uint16
	endPort       uint16
	host          string
	port          uint16
	addr          string
	size          atomic.Uint64
	current       atomic.Uint64
	bo            backoff.Backoff
	eg            errgroup.Group
	group         singleflight.Group[Conn]
	dopts         []DialOption
	dialTimeout   time.Duration
	roccd         time.Duration // reconnection old connection closing duration
	closing       atomic.Bool
	pmu           sync.RWMutex
	isIP          bool
	resolveDNS    bool
	reconnectHash atomic.Pointer[string]
}

const defaultPoolSize = 4

var (
	mu      sync.RWMutex
	metrics map[string]int64 = make(map[string]int64)
)

func New(ctx context.Context, opts ...Option) (c Conn, err error) {
	p := new(pool)

	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}

	p.init(true)
	p.closing.Store(false)

	p.group = singleflight.New[Conn]()

	var (
		isIPv4, isIPv6 bool
		port           uint16
	)
	p.host, p.port, _, isIPv4, isIPv6, err = net.Parse(p.addr)
	p.isIP = isIPv4 || isIPv6
	if err != nil {
		log.Warnf("failed to parse addr %s: %s", p.addr, err)
		if p.host == "" {
			var (
				ok      bool
				portStr string
			)
			p.host, portStr, ok = strings.Cut(p.addr, ":")
			if !ok {
				p.host = p.addr
			} else {
				portNum, err := strconv.ParseUint(portStr, 10, 16)
				if err != nil {
					p.port = uint16(portNum)
				}
			}
		}
		if p.port == 0 {
			port, err = p.scanGRPCPort(ctx)
			if err != nil {
				return nil, err
			}
			p.port = port
		}
		p.addr = net.JoinHostPort(p.host, p.port)
	}

	log.Debugf("initial connection will try dialing to %s, host: %s, port:%d, is IP Conn: %t", p.addr, p.host, p.port, p.isIP)
	conn, err := p.dial(ctx, p.addr)
	if err != nil {
		log.Warnf("grpc.New initial Dial check to %s returned error: %v", p.addr, err)
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Warn("failed to close connection:", err)
			}
		}

		port, err := p.scanGRPCPort(ctx)
		if err != nil {
			return nil, err
		}
		p.port = port
		p.addr = net.JoinHostPort(p.host, p.port)
		log.Debugf("fallback initial connection will try dialing to %s, host: %s, port:%d, is IP Conn: %t", p.addr, p.host, p.port, p.isIP)
		conn, err := p.dial(ctx, p.addr)
		if err != nil {
			if conn != nil {
				cerr := conn.Close()
				if cerr != nil {
					return nil, errors.Join(err, cerr)
				}
			}
			return nil, err
		}
	}
	if conn != nil {
		err = conn.Close()
		if err != nil {
			return nil, err
		}
	}

	if p.eg == nil {
		p.eg = errgroup.Get()
	}

	return p, nil
}

func (p *pool) init(force bool) {
	if p == nil {
		return
	}
	if p.Size() < 1 {
		p.size.Store(defaultPoolSize)
	}
	p.pmu.RLock()
	log.Debugf("initializing connection pool, addr: %s len: %d, cap: %d size: %d", p.addr, len(p.pool), cap(p.pool), p.Size())
	if force || p.pool == nil || cap(p.pool) == 0 || len(p.pool) == 0 {
		p.pmu.RUnlock()
		p.pmu.Lock()
		p.pool = make([]atomic.Pointer[poolConn], p.Size())
		p.pmu.Unlock()
	} else {
		p.pmu.RUnlock()
	}
}

func (p *pool) grow(size uint64) {
	if p == nil || p.Size() > size {
		return
	}
	l := p.Len()
	if l >= size {
		return
	}
	epool := make([]atomic.Pointer[poolConn], size-l) // expand pool
	log.Debugf("growing pool size %d o %d", l, size)
	p.pmu.Lock()
	if uint64(len(p.pool)) != l {
		epool = make([]atomic.Pointer[poolConn], size-uint64(len(p.pool))) // re-expand pool
	}
	p.pool = append(p.pool, epool...)
	p.pmu.Unlock()
	p.size.Store(size)
}

func (p *pool) load(idx int) (pc *poolConn) {
	if p == nil {
		return nil
	}
	p.pmu.RLock()
	if p.pool != nil && p.Size() > uint64(idx) && len(p.pool) > idx {
		pc = p.pool[idx].Load()
	}
	p.pmu.RUnlock()
	return pc
}

func (p *pool) store(idx int, pc *poolConn) {
	if p == nil {
		return
	}
	// p.init(false)
	p.pmu.Lock()
	if p.pool != nil && p.Size() > uint64(idx) && len(p.pool) > idx {
		p.pool[idx].Store(pc)
	}
	p.pmu.Unlock()
}

func (p *pool) loop(
	ctx context.Context, fn func(ctx context.Context, idx int, pc *poolConn) bool,
) (err error) {
	if p == nil || fn == nil {
		return nil
	}
	// p.init(false)
	p.pmu.RLock()
	defer p.pmu.RUnlock()
	var cnt int
	for idx, pool := range p.pool {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if p.Size() > uint64(idx) && len(p.pool) > idx {
				cnt++
				if !fn(ctx, idx, pool.Load()) {
					return nil
				}
			}
		}
	}
	if cnt == 0 {
		return errors.ErrGRPCPoolConnectionNotFound
	}
	return nil
}

func (p *pool) len() int {
	if p == nil {
		return 0
	}
	p.pmu.RLock()
	defer p.pmu.RUnlock()
	return len(p.pool)
}

func (p *pool) cap() int {
	if p == nil {
		return 0
	}
	p.pmu.RLock()
	defer p.pmu.RUnlock()
	return cap(p.pool)
}

func (p *pool) flush() {
	if p == nil {
		return
	}
	p.pmu.Lock()
	p.pool = nil
	p.pmu.Unlock()
}

func (p *pool) refreshConn(ctx context.Context, idx int, pc *poolConn, addr string) (err error) {
	if pc != nil && pc.addr == addr && isHealthy(ctx, pc.conn) {
		return nil
	}
	if pc != nil {
		log.Debugf("connection for %s pool %d/%d len %d is unhealthy trying to establish new pool member connection to %s", pc.addr, idx+1, p.Size(), p.Len(), addr)
	} else {
		log.Debugf("connection pool %d/%d len %d is empty, establish new pool member connection to %s", idx+1, p.Size(), p.Len(), addr)
	}
	conn, err := p.dial(ctx, addr)
	if err != nil {
		if pc != nil {
			if isHealthy(ctx, pc.conn) {
				log.Debugf("dialing new connection to %s failed,\terror: %v,\tbut existing connection to %s is healthy will keep existing connection", addr, err, pc.addr)
				return nil
			}
			if pc.conn != nil {
				p.eg.Go(safety.RecoverFunc(func() error {
					log.Debugf("waiting for invalid connection to %s to be closed...", pc.addr)
					err = pc.Close(ctx, p.roccd)
					if err != nil {
						log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, err)
					}
					return nil
				}))
			}
		}
		return errors.Join(err, errors.ErrInvalidGRPCClientConn(addr))
	}
	p.store(idx, &poolConn{
		conn: conn,
		addr: addr,
	})
	if pc != nil {
		p.eg.Go(safety.RecoverFunc(func() error {
			log.Debugf("waiting for old connection to %s to be closed...", pc.addr)
			err = pc.Close(ctx, p.roccd)
			if err != nil {
				log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, err)
			}
			return nil
		}))
	}
	return nil
}

func (p *pool) Connect(ctx context.Context) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}
	// p.init(false)

	log.Debugf("Connecting to addr: %s, host: %s, port: %d, isIP: %t, resolveDNS: %t", p.addr, p.host, p.port, p.isIP, p.resolveDNS)

	if p.isIP || !p.resolveDNS {
		return p.singleTargetConnect(ctx, p.addr)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		return p.singleTargetConnect(ctx, p.addr)
	}
	if len(ips) == 1 {
		return p.singleTargetConnect(ctx, ips[0])
	}
	return p.connect(ctx, ips...)
}

func (p *pool) connect(ctx context.Context, ips ...string) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	if uint64(len(ips)) > p.Size() {
		p.grow(uint64(len(ips)))
	}

	log.Debugf("connecting to ips: %v, port:  %d", ips, p.port)

	err = p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		addr := net.JoinHostPort(ips[idx%len(ips)], p.port)
		ierr := p.refreshConn(ctx, idx, pc, addr)
		if ierr != nil {
			if !errors.Is(ierr, context.DeadlineExceeded) &&
				!errors.Is(ierr, context.Canceled) {
				log.Warnf("An error occurred while dialing pool member connection to %s,\terror: %v", addr, ierr)
			} else {
				log.Debugf("Connect loop operation canceled while dialing pool member connection to %s,\terror: %v", addr, ierr)
				return false
			}
		}
		return true
	})
	if !errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		return p, err
	}

	hash := strings.Join(ips, "-")
	p.reconnectHash.Store(&hash)

	return p, nil
}

func (p *pool) singleTargetConnect(ctx context.Context, addr string) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	log.Debugf("connecting to single target addr: %s, host: %s, port: %d", p.addr, p.host, p.port)

	failCnt := 0
	err = p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		ierr := p.refreshConn(ctx, idx, pc, addr)
		if ierr != nil {
			if !errors.Is(ierr, context.DeadlineExceeded) &&
				!errors.Is(ierr, context.Canceled) {
				log.Warnf("An error occurred while dialing pool member connection to %s,\terror: %v", p.addr, ierr)
				failCnt++
				if p.isIP && (p.len() <= 2 || failCnt >= p.len()/3) {
					return false
				}
				return true
			} else {
				log.Debugf("Connect loop operation canceled while dialing pool member connection to %s,\terror: %v", p.addr, ierr)
				return false
			}
		}
		return true
	})
	if !errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		return p, err
	}
	p.reconnectHash.Store(&p.host)
	return p, nil
}

func (p *pool) Reconnect(ctx context.Context, force bool) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	log.Debug("Re-Connecting...")

	hash := p.reconnectHash.Load()
	if force || hash == nil || *hash == "" {
		log.Debug("re-connecting to all pool members")
		return p.Connect(ctx)
	}

	healthy := p.IsHealthy(ctx)
	if healthy {
		log.Debugf("dns connection is healthy, trying to reconnect using dns for more efficient ip balancing")
		if !p.isIP && p.resolveDNS && hash != nil && *hash != "" {
			ips, err := p.lookupIPAddr(ctx)
			if err != nil {
				return p, nil
			}
			if len(ips) == 1 {
				return p.singleTargetConnect(ctx, ips[0])
			}
			if *hash != strings.Join(ips, "-") {
				return p.connect(ctx, ips...)
			}
		}
		return p, nil
	}

	return p.Connect(ctx)
}

func (p *pool) Disconnect(ctx context.Context) (err error) {
	log.Debug("Disconnecting...")
	p.closing.Store(true)
	defer p.closing.Store(false)
	emap := make(map[string]error, p.len())
	err = p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		if pc != nil && pc.conn != nil {
			log.Debugf("closing connection for pool index: %d, addr: %s", idx, pc.addr)
			ierr := pc.conn.Close()
			if ierr != nil {
				if !errors.Is(ierr, context.DeadlineExceeded) &&
					!errors.Is(ierr, context.Canceled) {
					log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, ierr)
					emap[ierr.Error()] = err
				} else {
					log.Debugf("Disconnect loop operation canceled while closing pool member connection to %s,\terror: %v", pc.addr, ierr)
					return false
				}
			}
		}
		return true
	})
	p.flush()
	for _, e := range emap {
		err = errors.Join(err, e)
	}
	return err
}

func (p *pool) dial(ctx context.Context, addr string) (conn *ClientConn, err error) {
	do := func() (conn *ClientConn, err error) {
		ctx, cancel := context.WithTimeout(ctx, p.dialTimeout)
		defer cancel()
		log.Debugf("dialing to %s with timeout %s", addr, p.dialTimeout.String())
		conn, err = grpc.NewClient(addr, p.dopts...)
		if err != nil {
			if conn != nil {
				cerr := conn.Close()
				if cerr != nil {
					err = errors.Join(err, cerr)
				}
			}
			log.Debugf("failed to dial gRPC connection to %s,\terror: %v", addr, err)
			return nil, err
		}
		if !isHealthy(ctx, conn) {
			if conn != nil {
				err = conn.Close()
				if err != nil {
					err = errors.Join(errors.ErrGRPCClientConnNotFound(addr), err)
				} else {
					err = errors.ErrGRPCClientConnNotFound(addr)
				}
			}
			log.Debugf("connection for %s is unhealthy: %v", addr, err)
			return nil, err
		}
		return conn, nil
	}
	if p.bo != nil {
		retry := 0
		_, err = p.bo.Do(ctx, func(ctx context.Context) (r any, ret bool, err error) {
			log.Debugf("dialing to %s with backoff, retry: %d", addr, retry)
			conn, err = do()
			retry++
			return conn, err != nil, err
		})
		return conn, nil
	}

	return do()
}

func (p *pool) IsHealthy(ctx context.Context) (healthy bool) {
	if p == nil || p.closing.Load() {
		return false
	}
	log.Debug("Checking health...")
	var cnt, unhealthy int
	pl := p.len()
	err := p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		if pc == nil || !isHealthy(ctx, pc.conn) {
			if p.isIP {
				if pc != nil && pc.addr != "" {
					log.Debugf("unhealthy ip connection for pool index: %d, addr %s detected during health check, trying to re-connect", idx, pc.addr)
					err := p.refreshConn(ctx, idx, pc, pc.addr)
					if err != nil {
						// target addr cannot re-connect so, connection is unhealthy
						log.Error("failed to ip re-connecting to pool index: %d, addr %s", idx, pc.addr)
						unhealthy = pl - cnt
						return false
					}
					return true
				}
				log.Warn("unhealthy ip connection detected for pool index: %d, addr %s", idx, pc.addr)
				unhealthy = pl - cnt
				return false
			}
			addr := p.addr
			if pc != nil && pc.addr != "" && pc.addr != addr {
				addr = pc.addr
			}
			log.Warn("unhealthy dns connection for pool index: %d, addr %s detected during health check, trying to re-connect", idx, addr)
			// re-connect to last connected addr
			err := p.refreshConn(ctx, idx, pc, addr)
			if err != nil {
				if addr == p.addr {
					unhealthy++
					return true
				}
				// last connect addr is not dns and cannot connect then try dns
				err = p.refreshConn(ctx, idx, pc, p.addr)
				// dns addr cannot connect so, connection is unhealthy
				if err != nil {
					log.Error("failed to dns re-connecting to pool index: %d, addr %s", idx, pc.addr)
					unhealthy = pl - cnt
					return false
				}
			}
		}
		cnt++
		return true
	})
	if err != nil {
		log.Debugf("health check loop for addr=%s returned error: %v,\thealthy %d/%d", p.addr, err, pl-unhealthy, pl)
	}
	mu.Lock()
	metrics[p.addr] = int64(pl - unhealthy)
	mu.Unlock()

	if cnt == 0 {
		log.Debugf("no connection pool %d/%d found for %s,\thealthy %d/%d", cnt, pl, p.addr, pl-unhealthy, pl)
		return false
	}
	if p.isIP {
		// if ip pool connection, each connection target should be healthy
		return unhealthy == 0
	}

	// some pool target may unhealthy but pool client is healthy when unhealthy is less than pool length
	return unhealthy < pl
}

func (p *pool) Do(ctx context.Context, f func(conn *ClientConn) error) (err error) {
	if p == nil {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	idx, pc, ok := p.getHealthyConn(ctx, 0, p.Len())
	if !ok || pc == nil || pc.conn == nil {
		return errors.ErrGRPCClientConnNotFound(p.addr)
	}

	conn := pc.conn
	err = f(conn)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st != nil && st.Code() != codes.Canceled { // connection closing or closed
			if conn != nil {
				cerr := conn.Close()
				if cerr != nil {
					st, ok := status.FromError(cerr)
					if ok && st != nil && st.Code() != codes.Canceled { // connection closing or closed
						log.Warnf("Failed to close connection: %v", cerr)
					}
				}
			}
			rerr := p.refreshConn(ctx, idx, pc, p.addr)
			if rerr == nil {
				if newErr := f(p.load(idx).conn); newErr != nil {
					return errors.Join(err, newErr)
				}
				return nil
			}
			err = errors.Join(err, rerr)
		}
	}
	return err
}

func (p *pool) Get(ctx context.Context) (conn *ClientConn, ok bool) {
	_, pc, ok := p.getHealthyConn(ctx, 0, p.Len())
	if ok && pc != nil {
		return pc.conn, true
	}
	return nil, false
}

func (p *pool) getHealthyConn(
	ctx context.Context, cnt, retry uint64,
) (idx int, pc *poolConn, ok bool) {
	if p == nil || p.closing.Load() {
		return 0, nil, false
	}
	select {
	case <-ctx.Done():
		return 0, nil, false
	default:
	}
	pl := p.Len()
	if retry <= 0 || retry > math.MaxUint64-pl || pl <= 0 {
		if p.isIP {
			log.Warnf("failed to find gRPC IP connection pool for %s.\tlen(pool): %d,\tretried: %d,\tseems IP %s is unhealthy will going to disconnect...", p.addr, pl, cnt, p.addr)
			if err := p.Disconnect(ctx); err != nil {
				log.Debugf("failed to disconnect gRPC IP direct connection for %s,\terr: %v", p.addr, err)
			}
			return 0, nil, false
		}
		if pl > 0 {
			idx = int(p.current.Add(1) % pl)
		}
		if pc = p.load(idx); pc != nil && isHealthy(ctx, pc.conn) {
			return idx, pc, true
		}
		err := p.refreshConn(ctx, idx, pc, p.addr)
		if err == nil {
			return idx, p.load(idx), true
		}
		log.Warnf("failed to find gRPC connection pool for %s.\tlen(pool): %d,\tretried: %d,\terror: %v", p.addr, pl, cnt, err)
		return idx, nil, false
	}

	if pl > 0 {
		idx = int(p.current.Add(1) % pl)
		if pc = p.load(idx); pc != nil && isHealthy(ctx, pc.conn) {
			return idx, pc, true
		}
	}
	retry--
	cnt++
	return p.getHealthyConn(ctx, cnt, retry)
}

func (p *pool) Len() uint64 {
	return uint64(p.len())
}

func (p *pool) Size() uint64 {
	return p.size.Load()
}

func (p *pool) lookupIPAddr(ctx context.Context) (ips []string, err error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, p.host)
	if err != nil {
		log.Debugf("failed to lookup ip addr for %s \terr: %s", p.addr, err.Error())
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, errors.ErrGRPCLookupIPAddrNotFound(p.host)
	}

	ips = make([]string, 0, len(addrs))
	for _, ip := range addrs {
		ipStr := ip.String()
		var conn net.Conn
		addr := net.JoinHostPort(ipStr, p.port)
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
		conn, err := net.DialContext(ctx, net.TCP.String(), addr)
		cancel()
		if err != nil || conn == nil {
			log.Warnf("failed to initialize ping addr: %s,\terr: %s", addr, err.Error())
		} else {
			ips = append(ips, ipStr)
		}
		if conn != nil {
			err = conn.Close()
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Warn("failed to close connection:", err)
			}
		}
	}

	if len(ips) == 0 {
		return nil, errors.ErrGRPCLookupIPAddrNotFound(p.host)
	}

	slices.Sort(ips)

	return ips, nil
}

func (p *pool) scanGRPCPort(ctx context.Context) (port uint16, err error) {
	ports, err := net.ScanPorts(ctx, p.startPort, p.endPort, p.host)
	if err != nil {
		return 0, err
	}

	log.Debugf("starting to scan available gRPC ports from %v", ports)

	var conn *ClientConn
	for _, port := range ports {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			// try gRPC dialing to target port
			conn, err = grpc.NewClient(net.JoinHostPort(p.host, port), p.dopts...)
			if err == nil && isHealthy(ctx, conn) && conn.Close() == nil {
				// if no error and healthy the port is ready for gRPC
				log.Debugf("successuly found available gRPC port on %d", port)
				return port, nil
			}

			if conn != nil {
				_ = conn.Close()
			}
		}
	}
	return 0, errors.ErrInvalidGRPCPort(p.addr, p.host, p.port)
}

func (p *pool) IsIPConn() (isIP bool) {
	return p.isIP
}

func (p *pool) String() (str string) {
	if p == nil {
		return "<nil>"
	}
	var hash string
	rh := p.reconnectHash.Load()
	if rh != nil {
		hash = *rh
	}
	return fmt.Sprintf("addr: %s, host: %s, port %d, isIP: %t, resolveDNS: %t, hash: %s, pool_len: %d, pool_size: %d, current_seek: %d, dopt_len: %d, dial_timeout: %s, roccd: %s, closing: %t",
		p.addr,
		p.host,
		p.port,
		p.isIP,
		p.resolveDNS,
		hash,
		len(p.pool),
		p.size.Load(),
		p.current.Load(),
		len(p.dopts),
		p.dialTimeout.String(),
		p.roccd.String(),
		p.closing.Load())
}

func (pc *poolConn) Close(ctx context.Context, delay time.Duration) error {
	tdelay := delay / 10
	if tdelay < time.Millisecond*5 {
		tdelay = time.Millisecond * 5
	} else if tdelay > time.Minute {
		tdelay = time.Second * 5
	}
	tick := time.NewTicker(tdelay)
	defer tick.Stop()
	ctx, cancel := context.WithTimeout(ctx, delay)
	defer cancel()
	log.Debugf("closing all pool connection for %s with delay %s", pc.addr, delay.String())
	for {
		select {
		case <-ctx.Done():
			err := pc.conn.Close()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil && st.Code() != codes.Canceled { // connection closing or closed
					if ctx.Err() != nil &&
						!errors.Is(ctx.Err(), context.DeadlineExceeded) &&
						!errors.Is(ctx.Err(), context.Canceled) {
						return errors.Join(err, ctx.Err())
					}
					return err
				}
			}
			if ctx.Err() != nil &&
				!errors.Is(ctx.Err(), context.DeadlineExceeded) &&
				!errors.Is(ctx.Err(), context.Canceled) {
				return ctx.Err()
			}
			return nil
		case <-tick.C:
			switch pc.conn.GetState() {
			case connectivity.Idle, connectivity.Connecting, connectivity.TransientFailure:
				err := pc.conn.Close()
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil && st.Code() != codes.Canceled {
						return err
					}
				}
				return nil
			case connectivity.Shutdown:
				return nil
			}
		}
	}
}

func isHealthy(ctx context.Context, conn *ClientConn) bool {
	if conn == nil {
		log.Warn("gRPC target connection is nil")
		return false
	}
	state := conn.GetState()
	switch state {
	case connectivity.Ready:
		log.Debugf("gRPC target %s's connection status is Ready\tstatus: %s", conn.Target(), state.String())
		return true
	case connectivity.Connecting:
		log.Debugf("gRPC target %s's connection status will be Ready soon\tstatus: %s", conn.Target(), state.String())
		return true
	case connectivity.Idle:
		log.Debugf("gRPC target %s's connection status is Idle waiting for target\tstatus: %s\ttrying to re-connect...", conn.Target(), state.String())
		conn.Connect()
		if conn.WaitForStateChange(ctx, state) {
			return isHealthy(ctx, conn)
		}
		log.Errorf("gRPC target %s's connection status is not recovered\tstatus: %s", conn.Target(), state.String())
		return false
	case connectivity.Shutdown, connectivity.TransientFailure:
		log.Errorf("gRPC target %s's connection status is unhealthy\tstatus: %s", conn.Target(), state.String())
		return false
	}
	log.Errorf("gRPC target %s's connection status is unknown\tstatus: %s", conn.Target(), state.String())
	return false
}

func Metrics(context.Context) map[string]int64 {
	mu.RLock()
	defer mu.RUnlock()

	if len(metrics) == 0 {
		return nil
	}
	return maps.Clone(metrics)
}
