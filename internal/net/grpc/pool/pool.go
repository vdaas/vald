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

// Package pool provides grpc connection pool client
package pool

import (
	"context"
	"math"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	ClientConn = grpc.ClientConn
	DialOption = grpc.DialOption
)

type Conn interface {
	Connect(context.Context) (Conn, error)
	Disconnect() error
	Do(f func(*ClientConn) error) error
	Get() (*ClientConn, bool)
	IsHealthy(context.Context) bool
	Len() uint64
	Size() uint64
	Reconnect(ctx context.Context, force bool) (Conn, error)
}

type poolConn struct {
	conn *ClientConn
	addr string
}

type pool struct {
	pool          []atomic.Value
	startPort     uint16
	endPort       uint16
	host          string
	port          uint16
	addr          string
	size          uint64
	current       uint64
	bo            backoff.Backoff
	eg            errgroup.Group
	dopts         []DialOption
	dialTimeout   time.Duration
	roccd         time.Duration // reconnection old connection closing duration
	closing       atomic.Value
	isIP          bool
	resolveDNS    bool
	reconnectHash string
}

func New(ctx context.Context, opts ...Option) (c Conn, err error) {
	p := new(pool)

	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}

	if p.size < 1 {
		p.size = 1
	}

	p.pool = make([]atomic.Value, p.size)
	p.closing.Store(false)

	var isIPv4, isIPv6 bool
	p.host, p.port, _, isIPv4, isIPv6, err = net.Parse(p.addr)
	p.isIP = isIPv4 || isIPv6
	if err != nil {
		log.Warnf("failed to parse addr %s: %s", p.addr, err)
		if p.host == "" {
			var (
				ok   bool
				port string
			)
			p.host, port, ok = strings.Cut(p.addr, ":")
			if !ok {
				p.host = p.addr
			} else {
				portNum, err := strconv.ParseUint(port, 10, 16)
				if err != nil {
					p.port = uint16(portNum)
				}
			}
		}
		if p.port == 0 {
			err = p.scanGRPCPort(ctx)
			if err != nil {
				return nil, err
			}
		}
		p.addr = net.JoinHostPort(p.host, p.port)
	}

	conn, err := grpc.DialContext(ctx, p.addr, p.dopts...)
	if err != nil {
		log.Warn(err)
		err = p.scanGRPCPort(ctx)
		if err != nil {
			return nil, err
		}
		p.addr = net.JoinHostPort(p.host, p.port)
		conn, err = grpc.DialContext(ctx, p.addr, p.dopts...)
		if err != nil {
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

func (p *pool) Connect(ctx context.Context) (c Conn, err error) {
	if p == nil || p.closing.Load().(bool) {
		return p, nil
	}

	if p.pool == nil || cap(p.pool) == 0 || p.Len() == 0 {
		p.pool = make([]atomic.Value, p.size)
	}

	if p.isIP || !p.resolveDNS {
		return p.reconnectUnhealthy(ctx)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		return p.reconnectUnhealthy(ctx)
	}
	p.reconnectHash = strings.Join(ips, "-")

	for i := range p.pool {
		select {
		case <-ctx.Done():
			return p, nil
		default:
			var (
				conn   *ClientConn
				addr   = net.JoinHostPort(ips[i%len(ips)], p.port)
				pc, ok = p.load(i)
			)
			if ok && pc != nil && pc.addr == addr && isHealthy(pc.conn) {
				// TODO maybe we should check neighbour pool slice if new addrs come.
				continue
			}
			log.Debugf("establishing balanced connection to %s", addr)
			conn, err := p.dial(ctx, addr)
			if err != nil {
				log.Warnf("An error occurred during dialing to %s: %s", addr, err)
				continue
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: addr,
			})
			if pc != nil {
				p.eg.Go(safety.RecoverFunc(func() error {
					log.Debugf("waiting for old connection to %s to be closed...", pc.addr)
					err = pc.Close(ctx, p.roccd)
					if err != nil {
						log.Debugf("failed to close pool connection addr = %s\terror = %v", pc.addr, err)
					}
					return nil
				}))
			}
		}
	}

	return p, nil
}

func (p *pool) load(idx int) (pc *poolConn, ok bool) {
	if c := p.pool[idx].Load(); c != nil {
		pc, ok = c.(*poolConn)
	}
	return
}

func (p *pool) reconnectUnhealthy(ctx context.Context) (c Conn, err error) {
	p.reconnectHash = p.host
	failCnt := uint64(0)
	for i := range p.pool {
		select {
		case <-ctx.Done():
			return p, nil
		default:
			var (
				conn   *ClientConn
				pc, ok = p.load(i)
			)
			if ok && pc != nil && isHealthy(pc.conn) {
				continue
			}
			log.Debugf("establishing same connection to %s", p.addr)
			conn, err := p.dial(ctx, p.addr)
			if err != nil {
				failCnt++
				if p.isIP && (p.Len() <= 2 || failCnt >= p.Len()/3) {
					return nil, errors.ErrInvalidGRPCClientConn(p.addr)
				}
				continue
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: p.addr,
			})
			if pc != nil {
				p.eg.Go(safety.RecoverFunc(func() error {
					log.Debugf("waiting for old connection to %s to be closed...", pc.addr)
					err = pc.Close(ctx, p.roccd)
					if err != nil {
						log.Debugf("failed to close pool connection addr = %s\terror = %v", pc.addr, err)
					}
					return nil
				}))
			}
		}
	}
	return p, nil
}

func (p *pool) Disconnect() (err error) {
	for _, pool := range p.pool {
		pc, ok := pool.Load().(*poolConn)
		if ok && pc != nil && pc.conn != nil {
			err = pc.conn.Close()
			if err != nil {
				log.Debugf("failed to close pool connection addr = %s\terror = %v", pc.addr, err)
			}
		}
	}
	p.pool = nil
	return nil
}

func (p *pool) dial(ctx context.Context, addr string) (conn *ClientConn, err error) {
	if p.bo != nil {
		retry := 0
		_, err = p.bo.Do(ctx, func(ctx context.Context) (r interface{}, ret bool, err error) {
			log.Debugf("dialing to %s with backoff, retry: %d", addr, retry)
			ctx, cancel := context.WithTimeout(ctx, p.dialTimeout)
			defer cancel()
			conn, err = grpc.DialContext(ctx, addr, append(p.dopts, grpc.WithBlock())...)
			if err != nil {
				if conn != nil {
					err = errors.Wrap(conn.Close(), err.Error())
				}
				log.Debugf("failed to dial grpc connection to %s: %s", addr, err)
				retry++
				return nil, err != nil, err
			}
			if !isHealthy(conn) {
				if conn != nil {
					err = conn.Close()
				}
				if err != nil {
					err = errors.Wrapf(err, errors.ErrGRPCClientConnNotFound(addr).Error())
				} else {
					err = errors.ErrGRPCClientConnNotFound(addr)
				}
				log.Debugf("connection for %s is unhealthy: %s", addr, err)
				retry++
				return nil, err != nil, err
			}
			return conn, false, nil
		})
		return conn, nil
	}

	log.Debugf("dialing to %s", addr)
	ctx, cancel := context.WithTimeout(ctx, p.dialTimeout)
	defer cancel()
	conn, err = grpc.DialContext(ctx, addr, append(p.dopts, grpc.WithBlock())...)
	if err != nil || !isHealthy(conn) {
		log.Debugf("failed to dial pool connection addr = %s\terror = %v", addr, err)
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Debugf("failed to close pool connection addr = %s\terror = %v", addr, err)
			}
		}
		return nil, errors.Wrap(err, errors.ErrGRPCClientConnNotFound(addr).Error())
	}
	return conn, nil
}

func (p *pool) IsHealthy(ctx context.Context) bool {
	for i, pool := range p.pool {
		pc, ok := pool.Load().(*poolConn)
		if ok && pc != nil && pc.conn != nil && !isHealthy(pc.conn) {
			log.Debugf("connection for %s is unhealthy trying to dial for new connection", pc.addr)
			conn, err := p.dial(ctx, pc.addr)
			if err != nil {
				log.Warnf("failed to dial connection for %s", pc.addr)
				return false
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: pc.addr,
			})
			p.eg.Go(safety.RecoverFunc(func() error {
				log.Debugf("waiting for old connection to %s to be closed...", pc.addr)
				err = pc.Close(ctx, p.roccd)
				if err != nil {
					log.Warnf("failed to close old connection for %s,\terr: %v", pc.addr, err)
				}
				return nil
			}))
		}
	}
	return true
}

func (p *pool) Do(f func(conn *ClientConn) error) error {
	conn, ok := p.Get()
	if !ok {
		return errors.ErrGRPCClientConnNotFound(p.addr)
	}
	return f(conn)
}

func (p *pool) Get() (*ClientConn, bool) {
	return p.getHelthyConn(p.Len())
}

func (p *pool) getHelthyConn(retry uint64) (*ClientConn, bool) {
	if retry <= 0 || retry > math.MaxUint64-p.Len() || p.Len() <= 0 {
		log.Warnf("failed to find grpc pool connection for %s", p.addr)
		if p.isIP {
			log.Debugf("failure connection is IP connection trying to disconnect grpc connection for %s", p.addr)
			if err := p.Disconnect(); err != nil {
				log.Debugf("failed to disconnect grpc IP connection for %s,\terr: %v", p.addr, err)
			}
		}
		return nil, false
	}

	if res := p.pool[atomic.AddUint64(&p.current, 1)%p.Len()].Load(); res != nil {
		if pc, ok := res.(*poolConn); ok && pc != nil && isHealthy(pc.conn) {
			return pc.conn, true
		}
	}
	retry--
	return p.getHelthyConn(retry)
}

func (p *pool) Len() uint64 {
	return uint64(len(p.pool))
}

func (p *pool) Size() uint64 {
	return p.size
}

func (p *pool) lookupIPAddr(ctx context.Context) (ips []string, err error) {
	log.Debugf("resolving ip addr for %s", p.addr)
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, p.host)
	if err != nil {
		log.Debugf("failed to resolve ip addr for %s \terr: %s", p.addr, err.Error())
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
		if err != nil {
			log.Warnf("failed to initialize ping addr: %s,\terr: %s", addr, err.Error())
			continue
		}
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Warn("failed to close connection:", err)
			}
		}
		ips = append(ips, ipStr)
	}

	if len(ips) == 0 {
		return nil, errors.ErrGRPCLookupIPAddrNotFound(p.host)
	}

	sort.Strings(ips)

	return ips, nil
}

func (p *pool) Reconnect(ctx context.Context, force bool) (c Conn, err error) {
	healthy := p.IsHealthy(ctx)
	if p.isIP && p.reconnectHash != "" && !healthy {
		return nil, errors.ErrInvalidGRPCClientConn(p.addr)
	}

	if p.reconnectHash == "" {
		log.Debugf("connection history for %s not found starting first connection phase", p.addr)
		if p.isIP || !p.resolveDNS {
			return p.reconnectUnhealthy(ctx)
		}
		return p.Connect(ctx)
	}

	ips, err := p.lookupIPAddr(ctx)
	if err != nil || p.isIP {
		if !healthy {
			if p.isIP {
				return nil, errors.ErrInvalidGRPCClientConn(p.addr)
			}
			return p.reconnectUnhealthy(ctx)
		}
		return p, nil
	}
	if !healthy || p.reconnectHash != strings.Join(ips, "-") || force {
		return p.Connect(ctx)
	}
	return p, nil
}

func (p *pool) scanGRPCPort(ctx context.Context) (err error) {
	ports, err := net.ScanPorts(ctx, p.startPort, p.endPort, p.host)
	if err != nil {
		return err
	}
	for _, port := range ports {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if isGRPCPort(ctx, p.host, port) {
				p.port = port
				return nil
			}
		}
	}
	return errors.ErrInvalidGRPCPort(p.addr, p.host, p.port)
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
			if err != nil {
				if ctx.Err() != nil &&
					!errors.Is(ctx.Err(), context.DeadlineExceeded) &&
					!errors.Is(ctx.Err(), context.Canceled) {
					return errors.Wrap(err, ctx.Err().Error())
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
				if err != nil {
					return err
				}
				return nil
			case connectivity.Shutdown:
				return nil
			}
		}
	}
}

func isGRPCPort(ctx context.Context, host string, port uint16) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx,
		net.JoinHostPort(host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return false
	}
	return conn.Close() == nil
}

func isHealthy(conn *ClientConn) bool {
	if conn == nil {
		log.Warnf("grpc target %s's connection is nil", conn.Target())
		return false
	}
	state := conn.GetState()
	switch state {
	case connectivity.Ready:
		return true
	case connectivity.Idle, connectivity.Connecting:
		log.Debugf("grpc target %s's connection status will be Ready soon:\tstatus: %s", conn.Target(), state.String())
		return true
	case connectivity.Shutdown, connectivity.TransientFailure:
		log.Errorf("grpc target %s's connection status is unhealthy:\tstatus: %s", conn.Target(), state.String())
		return false
	}
	log.Errorf("grpc target %s's connection status is unknown:\tstatus: %s", conn.Target(), state.String())
	return false
}
