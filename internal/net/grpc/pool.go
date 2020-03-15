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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

type ClientConnPool interface {
	Connect(ctx context.Context) (ClientConnPool, error)
	Disconnect() error
	Do(f func(conn *ClientConn) error) error
	Get() (*ClientConn, bool)
	IsHealthy() bool
	Len() uint64
	Put(conn *ClientConn) error
	Reconnect(ctx context.Context, force bool) (ClientConnPool, error)
}

type clientConnPool struct {
	pool  atomic.Value
	host  string
	port  string
	isIP  bool
	size  uint64
	dopts []DialOption
}

type connPool struct {
	ctx     context.Context
	conn    *ClientConn // default connection
	group   singleflight.Group
	pool    sync.Pool
	addr    string
	host    string
	port    string
	size    uint64
	length  uint64
	dopts   []DialOption
	closing atomic.Value
}

func NewPool(ctx context.Context, addr string, size uint64, dopts ...DialOption) (ClientConnPool, error) {
	var ok bool
	host, port, isIP, err := net.Parse(addr)
	if err != nil {
		log.Warn(err)
		if len(host) == 0 {
			host = addr
		}
		p, ok := scanGRPCPort(ctx, host)
		if ok {
			port = p
		}
	}
	cp, err := newPool(ctx, host, port, size, dopts...)
	if err != nil {
		if cp != nil {
			err = errors.Wrap(cp.disconnect(), err.Error())
		}
		port, ok = scanGRPCPort(ctx, host)
		if !ok {
			return nil, err
		}
		cp, err = newPool(ctx, host, port, size, dopts...)
	}
	ccp := &clientConnPool{
		size:  size,
		dopts: dopts,
		host:  host,
		port:  port,
		isIP:  isIP,
	}
	ccp.pool.Store(cp)
	return ccp, nil
}

func scanGRPCPort(ctx context.Context, host string) (string, bool) {
	ports, err := net.ScanPorts(ctx, 1, 65535, host)
	if err != nil {
		log.Error(err)
		return "", false
	}
	for _, port := range ports {
		if isGRPCPort(ctx, host, port) {
			return strconv.Itoa(int(port)), true
		}
	}
	return "", false
}

func (c *clientConnPool) Connect(ctx context.Context) (ClientConnPool, error) {
	_, err := c.pool.Load().(*connPool).connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *clientConnPool) Disconnect() error {
	return c.pool.Load().(*connPool).disconnect()
}

func (c *clientConnPool) Do(f func(conn *ClientConn) error) error {
	return c.pool.Load().(*connPool).do(f)
}

func (c *clientConnPool) Get() (*ClientConn, bool) {
	return c.pool.Load().(*connPool).get()
}

func (c *clientConnPool) IsHealthy() bool {
	return c.pool.Load().(*connPool).isHealthy()
}

func (c *clientConnPool) Len() uint64 {
	return c.pool.Load().(*connPool).len()
}

func (c *clientConnPool) Put(conn *ClientConn) error {
	return c.pool.Load().(*connPool).put(conn)
}

func (c *clientConnPool) Reconnect(ctx context.Context, force bool) (ClientConnPool, error) {
	if !force && c.isIP {
		return c, nil
	}
	cp, err := newPool(ctx, c.host, c.port, c.size, c.dopts...)
	if err != nil {
		return nil, err
	}
	ocp := c.pool.Load().(*connPool)
	c.pool.Store(cp)
	return c, ocp.disconnect()
}

func newPool(ctx context.Context, host, port string, size uint64, dopts ...DialOption) (*connPool, error) {
	cp := &connPool{
		ctx:    ctx,
		addr:   host + ":" + port,
		host:   host,
		port:   port,
		size:   size,
		dopts:  dopts,
		length: 0,
	}
	cp.closing.Store(false)
	cp.pool = sync.Pool{
		New: func() interface{} {
			if cp.closing.Load().(bool) {
				return nil
			}
			if cp.conn != nil && isHealthy(cp.conn) {
				return cp.conn
			}
			ic, err, _ := cp.group.Do(cp.addr, func() (interface{}, error) {
				log.Warn("establishing new connection to " + cp.addr)
				conn, err := grpc.DialContext(ctx, cp.addr, cp.dopts...)
				if err != nil {
					log.Error(err)
					return nil, nil
				}
				if cp.conn != nil {
					cp.conn.Close()
				}
				cp.conn = conn
				return cp.conn, nil
			})
			if err != nil {
				if cp.conn != nil && isHealthy(cp.conn) {
					return cp.conn
				}
				log.Warn(err)
			}
			conn, ok := ic.(*ClientConn)
			if ok {
				return conn
			}
			return nil
		},
	}
	return cp.connect(ctx)
}

func (c *connPool) disconnect() (rerr error) {
	if c.closing.Load().(bool) {
		return nil
	}
	c.closing.Store(true)
	defer c.closing.Store(false)
	err := c.conn.Close()
	if err != nil {
		rerr = errors.Wrap(rerr, err.Error())
	}
	if c.size <= 1 {
		return nil
	}
	for i := uint64(0); i < uint64(math.Max(float64(atomic.LoadUint64(&c.length)), float64(c.size)))*2; i++ {
		conn, _ := c.get()
		if conn != nil {
			err = conn.Close()
			if err != nil {
				rerr = errors.Wrap(rerr, err.Error())
			}
		}
	}
	return
}

func (c *connPool) connect(ctx context.Context) (cp *connPool, err error) {
	if c.closing.Load().(bool) {
		return c, nil
	}

	if c.conn == nil || (c.conn != nil && !isHealthy(c.conn)) {
		conn, err := grpc.DialContext(ctx, c.addr, c.dopts...)
		if err != nil {
			log.Debugf("failed to dial pool connection addr = %s\terror = %v", c.addr, err)
			if conn != nil {
				return c, errors.Wrap(conn.Close(), err.Error())
			}
			return c, err
		}
		c.conn = conn
	}

	if c.size <= 1 {
		return c, nil
	}

	if atomic.LoadUint64(&c.length) > c.size {
		return c, nil
	}

	if net.IsLocal(c.host) {
		for atomic.LoadUint64(&c.length) > c.size {
			conn, err := grpc.DialContext(ctx, c.host+":"+c.port, c.dopts...)
			if err == nil {
				c.put(conn)
			} else {
				log.Debugf("failed to dial pool connection ip = %s\tport = %s\terror = %v", c.host, c.port, err)
				if conn != nil {
					return c, errors.Wrap(conn.Close(), err.Error())
				}
				return c, err
			}
		}
		return c, nil
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.host)
	if err != nil {
		for atomic.LoadUint64(&c.length) > c.size {
			conn, err := grpc.DialContext(ctx, c.addr, c.dopts...)
			if err == nil {
				c.put(conn)
			} else {
				log.Debugf("failed to dial pool connection addr = %s\terror = %v", c.addr, err)
				if conn != nil {
					return c, errors.Wrap(conn.Close(), err.Error())
				}
				return c, err
			}
		}
		return c, nil
	}

	if uint64(len(ips)) < c.size {
		for i := uint64(0); i < c.size/uint64(len(ips)); i++ {
			ips = append(ips, ips...)
		}
	}

	for _, ip := range ips {
		if atomic.LoadUint64(&c.length) > c.size {
			return c, nil
		}
		conn, err := grpc.DialContext(ctx, ip.String()+":"+c.port, c.dopts...)
		if err == nil {
			c.put(conn)
		} else {
			log.Debugf("failed to dial pool connection ip = %s\tport = %s\terror = %v", ip.String, c.port, err)
			if conn != nil {
				return c, errors.Wrap(conn.Close(), err.Error())
			}
			return c, err
		}
		if atomic.LoadUint64(&c.length) > c.size {
			return c, nil
		}
	}
	return c, nil
}

func (c *connPool) get() (*ClientConn, bool) {
	if c.size <= 1 {
		return c.conn, true
	}
	conn, ok := c.pool.Get().(*ClientConn)
	if !ok {
		return c.conn, true
	}
	atomic.AddUint64(&c.length, ^uint64(0))
	if conn == nil || !isHealthy(conn) {
		if c.closing.Load().(bool) {
			return nil, false
		}
		return c.conn, true
	}

	if conn == c.conn {
		return conn, true
	}

	return conn, false
}

func (c *connPool) put(conn *ClientConn) error {
	if c.size <= 1 {
		return nil
	}
	if conn != nil {
		if c.closing.Load().(bool) || atomic.LoadUint64(&c.length) > c.size {
			return conn.Close()
		}
		if conn == c.conn {
			return nil
		}
		atomic.AddUint64(&c.length, 1)
		c.pool.Put(conn)
	}
	return nil
}

func (c *connPool) do(f func(conn *ClientConn) error) (err error) {
	if c.size <= 1 {
		return f(c.conn)
	}
	conn, shared := c.get()
	if !shared {
		c.put(conn)
	}
	return f(conn)
}

func (c *connPool) isHealthy() bool {
	if c.conn == nil || !isHealthy(c.conn) {
		return false
	}
	if c.size <= 1 {
		return true
	}
	for i := uint64(0); i < c.size; i++ {
		conn, shared := c.get()
		if conn != nil && isHealthy(conn) {
			if !shared {
				c.put(conn)
			}
		} else {
			if conn != nil {
				conn.Close()
			}
			return false
		}
	}
	return true
}

func (c *connPool) len() uint64 {
	return atomic.LoadUint64(&c.length)
}

func isGRPCPort(ctx context.Context, host string, port uint16) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*3)
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
