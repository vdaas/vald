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
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

type ClientConnPool struct {
	ctx     context.Context
	pool    sync.Pool
	addr    string
	size    uint64
	length  uint64
	dopts   []DialOption
	closing atomic.Value
}

func NewPool(ctx context.Context, addr string, size uint64, dopts ...DialOption) (*ClientConnPool, error) {
	cp := &ClientConnPool{
		ctx:    ctx,
		addr:   addr,
		size:   size,
		dopts:  dopts,
		length: size,
	}
	cp.closing.Store(false)
	cp.pool = sync.Pool{
		New: func() interface{} {
			if cp.closing.Load().(bool) {
				return nil
			}
			conn, err := grpc.DialContext(ctx, addr, dopts...)
			if err != nil {
				log.Error(err)
				return nil
			}
			return conn
		},
	}
	return cp.Connect()
}

func (c *ClientConnPool) Disconnect() (rerr error) {
	if c.closing.Load().(bool) {
		return nil
	}
	c.closing.Store(true)
	defer c.closing.Store(false)
	for {
		conn := c.Get()
		if conn == nil {
			return
		}
		err := conn.Close()
		if err != nil {
			rerr = errors.Wrap(rerr, err.Error())
		}
	}
}

func (c *ClientConnPool) Connect() (cp *ClientConnPool, err error) {
	if c.closing.Load().(bool) {
		return nil, nil
	}
	for i := uint64(0); i < c.size; i++ {
		err = c.Do(func(conn *ClientConn) error {
			if isHealthy(conn) {
				return nil
			}
			return errors.ErrGRPCClientConnNotFound(c.addr)
		})
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *ClientConnPool) Get() *ClientConn {
	conn := c.pool.Get().(*ClientConn)
	atomic.AddUint64(&c.length, ^uint64(0))
	if conn == nil || !isHealthy(conn) {
		if c.closing.Load().(bool) {
			return nil
		}
		var err error
		conn, err = grpc.DialContext(c.ctx, c.addr, c.dopts...)
		if err != nil {
			log.Error(err)
			return nil
		}
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
		if isHealthy(conn) {
			c.Put(conn)
		} else {
			return false
		}
	}
	return true
}

func (c *ClientConnPool) Len() uint64 {
	return atomic.LoadUint64(&c.length)
}
