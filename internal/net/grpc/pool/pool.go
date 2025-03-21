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
	"math"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type (
	ClientConn = grpc.ClientConn
	DialOption = grpc.DialOption
)

type Conn interface {
	Connect(context.Context) (Conn, error)
	Disconnect() error
	Do(ctx context.Context, f func(*ClientConn) error) error
	Get(ctx context.Context) (conn *ClientConn, ok bool)
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

func New(ctx context.Context, opts ...Option) (c Conn, err error) {
	p := new(pool)

	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}

	p.init(true)
	p.closing.Store(false)

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

	conn, err := grpc.NewClient(p.addr, p.dopts...)
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
		conn, err = grpc.NewClient(p.addr, p.dopts...)
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
	p.init(false)
	p.pmu.RLock()
	if p.pool != nil && p.Size() > uint64(idx) && len(p.pool) > idx {
		p.pool[idx].Store(pc)
	}
	p.pmu.RUnlock()
}

func (p *pool) loop(
	ctx context.Context, fn func(ctx context.Context, idx int, pc *poolConn) bool,
) (err error) {
	if p == nil || fn == nil {
		return nil
	}
	p.init(false)
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
		log.Debugf("connection for %s pool %d/%d is unhealthy trying to establish new pool member connection to %s", pc.addr, idx+1, p.Size(), addr)
	} else {
		log.Debugf("connection pool %d/%d is empty, establish new pool member connection to %s", idx+1, p.Size(), addr)
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

	p.init(false)

	if p.isIP || !p.resolveDNS {
		return p.singleTargetConnect(ctx)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		return p.singleTargetConnect(ctx)
	}
	return p.connect(ctx, ips...)
}

func (p *pool) connect(ctx context.Context, ips ...string) (c Conn, err error) {
	if uint64(len(ips)) > p.Size() {
		p.grow(uint64(len(ips)))
	}

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

func (p *pool) Reconnect(ctx context.Context, force bool) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	hash := p.reconnectHash.Load()
	if force || hash == nil || *hash == "" {
		return p.Connect(ctx)
	}

	healthy := p.IsHealthy(ctx)
	if healthy {
		if !p.isIP && p.resolveDNS && hash != nil && *hash != "" {
			ips, err := p.lookupIPAddr(ctx)
			if err != nil {
				return p, nil
			}
			if *hash != strings.Join(ips, "-") {
				return p.connect(ctx, ips...)
			}
		}
		return p, nil
	}

	return p.Connect(ctx)
}

func (p *pool) singleTargetConnect(ctx context.Context) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	failCnt := 0
	err = p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		ierr := p.refreshConn(ctx, idx, pc, p.addr)
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

func (p *pool) Disconnect() (err error) {
	ctx := context.Background()
	p.closing.Store(true)
	defer p.closing.Store(false)
	emap := make(map[string]error, p.len())
	err = p.loop(ctx, func(ctx context.Context, _ int, pc *poolConn) bool {
		if pc != nil && pc.conn != nil {
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

	log.Debugf("dialing to %s", addr)
	return do()
}

func (p *pool) IsHealthy(ctx context.Context) (healthy bool) {
	if p == nil || p.closing.Load() {
		return false
	}
	var cnt, unhealthy int
	pl := p.len()
	err := p.loop(ctx, func(ctx context.Context, idx int, pc *poolConn) bool {
		if pc == nil || !isHealthy(ctx, pc.conn) {
			if p.isIP {
				if pc != nil && pc.addr != "" {
					err := p.refreshConn(ctx, idx, pc, pc.addr)
					if err != nil {
						// target addr cannot re-connect so, connection is unhealthy
						unhealthy++
						return false
					}
					return true
				}
				return false
			}
			addr := p.addr
			if pc != nil {
				addr = pc.addr
			}
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
	idx, conn, ok := p.getHealthyConn(ctx, 0, p.Len())
	if !ok || conn == nil {
		return errors.ErrGRPCClientConnNotFound(p.addr)
	}
	err = f(conn)
	if errors.Is(err, grpc.ErrClientConnClosing) {
		if conn != nil {
			if cerr := conn.Close(); cerr != nil && !errors.Is(cerr, grpc.ErrClientConnClosing) {
				log.Warnf("Failed to close connection: %v", cerr)
			}
		}
		conn, err = p.dial(ctx, p.addr)
		if err == nil && conn != nil && isHealthy(ctx, conn) {
			p.store(idx, &poolConn{
				conn: conn,
				addr: p.addr,
			})
			if newErr := f(conn); newErr != nil {
				return errors.Join(err, newErr)
			}
			return nil
		}
	}
	return err
}

func (p *pool) Get(ctx context.Context) (conn *ClientConn, ok bool) {
	_, conn, ok = p.getHealthyConn(ctx, 0, p.Len())
	return conn, ok
}

func (p *pool) getHealthyConn(
	ctx context.Context, cnt, retry uint64,
) (idx int, conn *ClientConn, ok bool) {
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
			if err := p.Disconnect(); err != nil {
				log.Debugf("failed to disconnect gRPC IP direct connection for %s,\terr: %v", p.addr, err)
			}
			return 0, nil, false
		}
		if pl > 0 {
			idx = int(p.current.Add(1) % pl)
		}
		if pc := p.load(idx); pc != nil && isHealthy(ctx, pc.conn) {
			return idx, pc.conn, true
		}
		conn, err := p.dial(ctx, p.addr)
		if err == nil && conn != nil && isHealthy(ctx, conn) {
			p.store(idx, &poolConn{
				conn: conn,
				addr: p.addr,
			})
			return idx, conn, true
		}
		log.Warnf("failed to find gRPC connection pool for %s.\tlen(pool): %d,\tretried: %d,\terror: %v", p.addr, pl, cnt, err)
		return idx, nil, false
	}

	if pl > 0 {
		idx = int(p.current.Add(1) % pl)
		if pc := p.load(idx); pc != nil && isHealthy(ctx, pc.conn) {
			return idx, pc.conn, true
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
	return fmt.Sprintf("addr: %s, host: %s, port %d, isIP: %v, resolveDNS: %v, hash: %s, pool_size: %d, current_seek: %d, dopt_len: %d, dial_timeout: %v, roccd: %v, closing: %v",
		p.addr,
		p.host,
		p.port,
		p.isIP,
		p.resolveDNS,
		hash,
		p.size.Load(),
		p.current.Load(),
		len(p.dopts),
		p.dialTimeout,
		p.roccd,
		p.closing.Load())
}

func (pc *poolConn) Close(ctx context.Context, delay time.Duration) error {
	tdelay := delay / 10
	if tdelay < time.Millisecond*200 {
		tdelay = time.Millisecond * 200
	} else if tdelay > time.Minute {
		tdelay = time.Second * 5
	}
	tick := time.NewTicker(tdelay)
	defer tick.Stop()
	ctx, cancel := context.WithTimeout(ctx, delay)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			err := pc.conn.Close()
			if err != nil && !errors.Is(err, grpc.ErrClientConnClosing) {
				if ctx.Err() != nil &&
					!errors.Is(ctx.Err(), context.DeadlineExceeded) &&
					!errors.Is(ctx.Err(), context.Canceled) {
					return errors.Join(err, ctx.Err())
				}
				return err
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
				if err != nil && !errors.Is(err, grpc.ErrClientConnClosing) {
					return err
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
		return true
	case connectivity.Connecting:
		log.Debugf("gRPC target %s's connection status will be Ready soon\tstatus: %s", conn.Target(), state.String())
		return true
	case connectivity.Idle:
		log.Debugf("gRPC target %s's connection status is waiting for target\tstatus: %s\ttrying to re-connect...", conn.Target(), state.String())
		conn.Connect()
		if conn.WaitForStateChange(ctx, state) {
			state = conn.GetState()
			if state == connectivity.Ready || state == connectivity.Connecting {
				log.Debugf("gRPC target %s's connection status enabled for target\tstatus: %s", conn.Target(), state.String())
				return true
			}
		}
		return false
	case connectivity.Shutdown, connectivity.TransientFailure:
		log.Errorf("gRPC target %s's connection status is unhealthy\tstatus: %s", conn.Target(), state.String())
		return false
	}
	log.Errorf("gRPC target %s's connection status is unknown\tstatus: %s", conn.Target(), state.String())
	return false
}
