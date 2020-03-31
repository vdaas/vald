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
	"math"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

type ClientConnPool struct {
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

const (
	localIPv4   = "127.0.0.1"
	localHost   = "localhost"
	defaultPort = "80"
)

func NewPool(ctx context.Context, addr string, size uint64, dopts ...DialOption) (*ClientConnPool, error) {
	if !strings.HasPrefix(addr, "::") && strings.HasPrefix(addr, ":") {
		addr = localIPv4 + addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
		port = defaultPort
	}
	cp := &ClientConnPool{
		ctx:    ctx,
		addr:   addr,
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
	return cp.Connect(ctx)
}

func (c *ClientConnPool) Disconnect() (rerr error) {
	if c.closing.Load().(bool) {
		return nil
	}
	c.closing.Store(true)
	defer c.closing.Store(false)
	err := c.conn.Close()
	if err != nil {
		rerr = errors.Wrap(rerr, err.Error())
	}
	for i := uint64(0); i < uint64(math.Max(float64(atomic.LoadUint64(&c.length)), float64(c.size)))*2; i++ {
		conn, _ := c.Get()
		if conn != nil {
			err = conn.Close()
			if err != nil {
				rerr = errors.Wrap(rerr, err.Error())
			}
		}
	}
	return
}

func (c *ClientConnPool) Connect(ctx context.Context) (cp *ClientConnPool, err error) {
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
			return nil, err
		}
		c.conn = conn
	}

	if atomic.LoadUint64(&c.length) > c.size {
		return c, nil
	}

	if c.host == localHost ||
		c.host == localIPv4 {
		for atomic.LoadUint64(&c.length) > c.size {
			conn, err := grpc.DialContext(ctx, localIPv4+":"+c.port, c.dopts...)
			if err == nil {
				c.Put(conn)
			} else {
				log.Debugf("failed to dial pool connection ip = %s\tport = %s\terror = %v", localIPv4, c.port, err)
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
				c.Put(conn)
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
			c.Put(conn)
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

func (c *ClientConnPool) Get() (*ClientConn, bool) {
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

func (c *ClientConnPool) Put(conn *ClientConn) error {
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

func (c *ClientConnPool) Do(f func(conn *ClientConn) error) (err error) {
	conn, shared := c.Get()
	if !shared {
		c.Put(conn)
	}
	return f(conn)
}

func (c *ClientConnPool) IsHealthy() bool {
	for i := uint64(0); i < c.size; i++ {
		conn, shared := c.Get()
		if conn != nil && isHealthy(conn) {
			if !shared {
				c.Put(conn)
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

func (c *ClientConnPool) Len() uint64 {
	return atomic.LoadUint64(&c.length)
}
