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
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

type ClientConnPool struct {
	ctx     context.Context
	conn    *ClientConn // default connection
	pool    sync.Pool
	addr    string
	host    string
	port    string
	size    uint64
	length  uint64
	dopts   []DialOption
	closing atomic.Value
}

func NewPool(ctx context.Context, addr string, size uint64, dopts ...DialOption) (*ClientConnPool, error) {
	if !strings.HasPrefix(addr, "::") && strings.HasPrefix(addr, ":") {
		addr = "127.0.0.1" + addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
		port = "80"
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
			conn, err := grpc.DialContext(ctx, cp.addr, cp.dopts...)
			if err != nil {
				log.Error(err)
				return nil
			}
			return conn
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
	for {
		conn := c.Get()
		if conn == nil {
			return
		}
		err = conn.Close()
		if err != nil {
			rerr = errors.Wrap(rerr, err.Error())
		}
	}
}

func (c *ClientConnPool) Connect(ctx context.Context) (cp *ClientConnPool, err error) {
	if c.closing.Load().(bool) {
		return nil, nil
	}
	if c.conn == nil || (c.conn != nil && !isHealthy(c.conn)) {
		conn, err := grpc.DialContext(ctx, c.addr, c.dopts...)
		if err == nil {
			c.conn = conn
		}
	}

	if c.host == "localhost" ||
		c.host == "127.0.0.1" {
		for {
			conn, err := grpc.DialContext(ctx, "127.0.0.1:"+c.port, c.dopts...)
			if err == nil {
				c.Put(conn)
			}
			if atomic.LoadUint64(&c.length) > c.size {
				return c, nil
			}
		}
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.host)
	if err != nil {
		for {
			conn, err := grpc.DialContext(ctx, c.addr, c.dopts...)
			if err == nil {
				c.Put(conn)
			}
			if atomic.LoadUint64(&c.length) > c.size {
				return c, nil
			}
		}
	}

	if uint64(len(ips)) < c.size {
		for i := uint64(0); i < c.size/uint64(len(ips)); i++ {
			ips = append(ips, ips...)
		}
	}

	for _, ip := range ips {
		conn, err := grpc.DialContext(ctx, ip.String()+":"+c.port, c.dopts...)
		if err == nil {
			c.Put(conn)
		}
		if atomic.LoadUint64(&c.length) > c.size {
			return c, nil
		}
	}
	return c, nil
}

func (c *ClientConnPool) Get() *ClientConn {
	conn, ok := c.pool.Get().(*ClientConn)
	if !ok {
		return c.conn
	}
	atomic.AddUint64(&c.length, ^uint64(0))
	if conn == nil || !isHealthy(conn) {
		if c.closing.Load().(bool) {
			return nil
		}
		return c.conn
	}
	return conn
}

func (c *ClientConnPool) Put(conn *ClientConn) error {
	if conn != nil {
		if c.closing.Load().(bool) {
			return nil
		}
		if atomic.LoadUint64(&c.length) > c.size {
			return conn.Close()
		}
		atomic.AddUint64(&c.length, 1)
		c.pool.Put(conn)
	}
	return nil
}

func (c *ClientConnPool) Do(f func(conn *ClientConn) error) (err error) {
	conn := c.Get()
	err = f(conn)
	c.Put(conn)
	return err
}

func (c *ClientConnPool) Healthy() bool {
	for i := uint64(0); i < c.size; i++ {
		conn := c.Get()
		if conn != nil && isHealthy(conn) {
			c.Put(conn)
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
