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

// Package pool provides grpc connection pool client
package pool

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientConn = grpc.ClientConn
type DialOption = grpc.DialOption

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
	dopts         []DialOption
	closing       atomic.Value
	isIP          bool
	reconnectHash string
}

func New(ctx context.Context, opts ...Option) (c Conn, err error) {
	p := new(pool)

	for _, opt := range append(defaultOpts, opts...) {
		opt(p)
	}

	if p.size < 1 {
		p.size = 1
	}

	p.pool = make([]atomic.Value, p.size)
	p.closing.Store(false)

	p.host, p.port, p.isIP, err = net.Parse(p.addr)
	if err != nil {
		log.Warn(err)
		if len(p.host) == 0 {
			p.host = strings.SplitN(p.addr, ":", 2)[0]
		}
		err = p.scanGRPCPort(ctx)
		if err != nil {
			return nil, err
		}
		p.addr = fmt.Sprintf("%s:%d", p.host, p.port)
	}

	conn, err := grpc.DialContext(ctx, p.addr, p.dopts...)
	if err != nil {
		err = p.scanGRPCPort(ctx)
		if err != nil {
			return nil, err
		}
		p.addr = fmt.Sprintf("%s:%d", p.host, p.port)
	}
	if conn != nil {
		err = conn.Close()
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *pool) Connect(ctx context.Context) (c Conn, err error) {
	if p.closing.Load().(bool) {
		return p, nil
	}

	if p.pool == nil || cap(p.pool) == 0 || len(p.pool) == 0 {
		p.pool = make([]atomic.Value, p.size)
	}

	if p.isIP {
		return p.connect(ctx)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		return p.connect(ctx)
	}
	p.reconnectHash = strings.Join(ips, "-")

	for i := range p.pool {
		select {
		case <-ctx.Done():
			return p, nil
		default:
			var (
				c    = p.pool[i].Load()
				conn *ClientConn
				addr = fmt.Sprintf("%s:%d", ips[i%len(ips)], p.port)
			)
			log.Debugf("establishing balanced connection to %s", addr)
			conn, err := p.dial(ctx, addr)
			if err != nil {
				continue
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: addr,
			})
			if c != nil {
				pc, ok := c.(*poolConn)
				if ok && pc != nil && pc.conn != nil {
					err = pc.conn.Close()
					if err != nil {
						log.Debugf("failed to close pool connection addr = %s\terror = %v", pc.addr, err)
					}
				}
			}
		}
	}

	return p, nil
}

func (p *pool) connect(ctx context.Context) (c Conn, err error) {
	p.reconnectHash = p.host
	for i := range p.pool {
		select {
		case <-ctx.Done():
			return p, nil
		default:
			var (
				c    = p.pool[i].Load()
				conn *ClientConn
			)
			conn, err := p.dial(ctx, p.addr)
			if err != nil {
				continue
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: p.addr,
			})
			if c != nil {
				pc, ok := c.(*poolConn)
				if ok && pc != nil && pc.conn != nil {
					err = pc.conn.Close()
					if err != nil {
						log.Debugf("failed to close pool connection addr = %s\terror = %v", pc.addr, err)
					}
				}
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
		var res interface{}
		res, err = p.bo.Do(ctx, func() (interface{}, error) {
			conn, err := grpc.DialContext(ctx, addr, p.dopts...)
			if err != nil {
				if conn != nil {
					err = errors.Wrap(conn.Close(), err.Error())
				}
				return nil, err
			}
			if !isHealthy(conn) {
				return nil, errors.ErrGRPCClientConnNotFound(addr)
			}
			return conn, nil
		})
		var ok bool
		conn, ok = res.(*ClientConn)
		if !ok {
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
	} else {
		conn, err = grpc.DialContext(ctx, addr, p.dopts...)
	}
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
			conn, err := p.dial(ctx, pc.addr)
			if err != nil {
				return false
			}
			p.pool[i].Store(&poolConn{
				conn: conn,
				addr: pc.addr,
			})
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
	return p.get(0)
}

func (p *pool) get(retry int) (*ClientConn, bool) {
	if atomic.LoadUint64(&p.current) >= math.MaxUint64-2 {
		atomic.StoreUint64(&p.current, 0)
	}
	res := p.pool[atomic.AddUint64(&p.current, 1)%uint64(len(p.pool))].Load()
	if res != nil {
		pc, ok := res.(*poolConn)
		if ok && isHealthy(pc.conn) {
			return pc.conn, true
		}
	}
	if retry > len(p.pool) {
		return nil, false
	}
	retry++
	return p.get(retry)
}

func (p *pool) Len() uint64 {
	return uint64(len(p.pool))
}

func (p *pool) Size() uint64 {
	return p.size
}

func (p *pool) lookupIPAddr(ctx context.Context) (ips []string, err error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, p.host)
	if err != nil {
		return nil, err
	}
	ips = make([]string, 0, len(addrs))

	const network = "tcp"
	for _, ip := range addrs {
		ipStr := ip.String()
		if net.IsIPv6(ipStr) && !strings.Contains(ipStr, "[") {
			ipStr = fmt.Sprintf("[%s]", ipStr)
		}
		var conn net.Conn
		addr := fmt.Sprintf("%s:%d", ipStr, p.port)
		if net.DefaultResolver.Dial != nil {
			ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
			conn, err = net.DefaultResolver.Dial(ctx, network, addr)
			cancel()
		} else {
			conn, err = net.Dial(network, addr)
		}
		if err != nil {
			log.Warn(err)
			continue
		}
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Warn(err)
			}
		}
		ips = append(ips, ipStr)
	}

	sort.Strings(ips)

	return ips, nil
}

func (p *pool) Reconnect(ctx context.Context, force bool) (c Conn, err error) {
	if len(p.reconnectHash) == 0 {
		return p.Connect(ctx)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		if !p.IsHealthy(ctx) {
			return p.connect(ctx)
		}
		return p, nil
	}
	if !p.IsHealthy(ctx) || p.reconnectHash != strings.Join(ips, "-") || force {
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

func isGRPCPort(ctx context.Context, host string, port uint16) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*5)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

func isHealthy(conn *ClientConn) bool {
	return conn != nil &&
		conn.GetState() != connectivity.Shutdown &&
		conn.GetState() != connectivity.TransientFailure
}
