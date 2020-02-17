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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Server = grpc.Server
type ClientConn = grpc.ClientConn
type CallOption = grpc.CallOption
type DialOption = grpc.DialOption

type Client interface {
	StartConnectionMonitor(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, addr string, dopts ...DialOption) error
	Disconnect(addr string) error
	Range(ctx context.Context,
		f func(addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	RangeConcurrent(ctx context.Context,
		concurrency int,
		f func(addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	OrderedRange(ctx context.Context,
		order []string,
		f func(addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	OrderedRangeConcurrent(ctx context.Context,
		order []string,
		concurrency int,
		f func(addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	Do(ctx context.Context,
		addr string, f func(conn *ClientConn,
			copts ...CallOption) (interface{}, error)) (interface{}, error)
	GetAddrs() ([]string, []string)
	GetDialOption() []DialOption
	GetCallOption() []CallOption
	Close() error
}

type gRPCClient struct {
	addrs       []string
	poolSize    uint64
	clientCount uint64
	conns       grpcConns
	hcDur       time.Duration
	dopts       []DialOption
	copts       []CallOption
	eg          errgroup.Group
	bo          backoff.Backoff
}

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
		err = c.Do(func(conn *grpc.ClientConn) error {
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

func (c *ClientConnPool) Get() *grpc.ClientConn {
	conn := c.pool.Get().(*grpc.ClientConn)
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

func (c *ClientConnPool) Put(conn *grpc.ClientConn) error {
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

func (c *ClientConnPool) Do(f func(conn *grpc.ClientConn) error) (err error) {
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

func New(opts ...Option) (c Client) {
	g := new(gRPCClient)

	for _, opt := range append(defaultOpts, opts...) {
		opt(g)
	}

	return g
}

func (g *gRPCClient) StartConnectionMonitor(ctx context.Context) (<-chan error, error) {
	if g.addrs == nil || len(g.addrs) == 0 {
		return nil, errors.ErrGRPCTargetAddrNotFound
	}

	ech := make(chan error, len(g.addrs))

	for _, addr := range g.addrs {
		if len(addr) != 0 {
			err := g.Connect(ctx, addr, grpc.WithBlock())
			if err != nil {
				log.Error(err)
				ech <- err
			}
		}
	}

	if len(g.addrs) != 0 && atomic.LoadUint64(&g.clientCount) == 0 {
		return nil, errors.ErrGRPCClientConnNotFound(strings.Join(g.addrs, ",\t"))
	}

	g.eg.Go(safety.RecoverFunc(func() (err error) {
		tick := time.NewTicker(g.hcDur)
		defer close(ech)
		defer g.Close()
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tick.C:
				reconnList := make([]string, 0, int(atomic.LoadUint64(&g.clientCount)))
				g.conns.Range(func(addr string, pool *ClientConnPool) bool {
					if len(addr) != 0 && !pool.Healthy() {
						reconnList = append(reconnList, addr)
					}
					return true
				})

				for _, addr := range reconnList {
					if g.bo != nil {
						_, err = g.bo.Do(ctx, func() (interface{}, error) {
							return nil, g.Connect(ctx, addr, g.dopts...)
						})
					} else {
						err = g.Connect(ctx, addr, g.dopts...)
					}
					if err != nil {
						log.Error(err)
						ech <- err
					}
				}
			}
		}
	}))
	return ech, nil
}

func (g *gRPCClient) Range(ctx context.Context,
	f func(addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	wrapf := func(addr string, conn *ClientConn, copts ...CallOption) (err error) {
		return f(addr, conn, copts...)
	}
	g.conns.Range(func(addr string, pool *ClientConnPool) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			err := pool.Do(func(conn *ClientConn) (err error) {
				if g.bo != nil {
					_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
						err = wrapf(addr, conn, g.copts...)
						return
					})
				} else {
					err = wrapf(addr, conn, g.copts...)
				}
				return nil
			})
			if err != nil {
				rerr = errors.Wrap(rerr, errors.ErrRPCCallFailed(addr, err).Error())
			}
		}
		return true
	})
	return rerr
}

func (g *gRPCClient) RangeConcurrent(ctx context.Context,
	concurrency int,
	f func(addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	wrapf := func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) (err error) {
		return f(addr, conn, copts...)
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	g.conns.Range(func(addr string, pool *ClientConnPool) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			select {
			case <-egctx.Done():
				return nil
			default:
				return pool.Do(func(conn *ClientConn) (err error) {
					if g.bo != nil {
						_, err = g.bo.Do(egctx, func() (r interface{}, err error) {
							err = wrapf(egctx, addr, conn, g.copts...)
							return
						})
					} else {
						err = wrapf(egctx, addr, conn, g.copts...)
					}
					if err != nil {
						return errors.Wrap(rerr, errors.ErrRPCCallFailed(addr, err).Error())
					}
					return nil
				})
			}
		}))
		return true
	})
	return eg.Wait()
}

func (g *gRPCClient) OrderedRange(ctx context.Context,
	orders []string,
	f func(addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	wrapf := func(addr string, conn *ClientConn, copts ...CallOption) (err error) {
		return f(addr, conn, copts...)
	}
	var err error
	for _, addr := range orders {
		select {
		case <-ctx.Done():
			return nil
		default:
			pool, ok := g.conns.Load(addr)
			if ok {
				if err = pool.Do(func(conn *ClientConn) (err error) {
					if g.bo != nil {
						_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
							err = wrapf(addr, conn, g.copts...)
							return
						})
						return err
					}
					return wrapf(addr, conn, g.copts...)
				}); err != nil {
					rerr = errors.Wrap(rerr, errors.ErrRPCCallFailed(addr, err).Error())
				}
			}
		}
	}
	return rerr
}

func (g *gRPCClient) OrderedRangeConcurrent(ctx context.Context,
	orders []string,
	concurrency int,
	f func(addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	wrapf := func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) (err error) {
		return f(addr, conn, copts...)
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	for _, addr := range orders {
		pool, ok := g.conns.Load(addr)
		if ok {
			eg.Go(safety.RecoverFunc(func() (err error) {
				select {
				case <-ctx.Done():
					return nil
				default:
					return pool.Do(func(conn *ClientConn) (err error) {
						if g.bo != nil {
							_, err = g.bo.Do(egctx, func() (r interface{}, err error) {
								err = wrapf(egctx, addr, conn, g.copts...)
								return
							})
						} else {
							err = wrapf(egctx, addr, conn, g.copts...)
						}

						if err != nil {
							return errors.ErrRPCCallFailed(addr, err)
						}
						return nil
					})
				}
			}))
		}
	}
	return eg.Wait()
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(conn *ClientConn,
		copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	wrapf := func(conn *ClientConn, copts ...CallOption) (ret interface{}, err error) {
		return f(conn, copts...)
	}
	pool, ok := g.conns.Load(addr)
	if !ok {
		g.Connect(ctx, addr)
		pool, ok = g.conns.Load(addr)
		if !ok {
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
	}
	err = pool.Do(func(conn *ClientConn) error {
		if g.bo != nil {
			data, err = g.bo.Do(ctx, func() (r interface{}, err error) {
				r, err = wrapf(conn, g.copts...)
				if err != nil {
					return nil, err
				}
				return r, nil
			})
		} else {
			data, err = wrapf(conn, g.copts...)
		}
		return err
	})
	if err != nil {
		return nil, errors.ErrRPCCallFailed(addr, err)
	}
	return
}

func (g *gRPCClient) GetDialOption() []DialOption {
	return g.dopts
}

func (g *gRPCClient) GetCallOption() []CallOption {
	return g.copts
}

func (g *gRPCClient) Connect(ctx context.Context, addr string, dopts ...DialOption) (err error) {
	pool, ok := g.conns.Load(addr)
	if ok {
		if pool.Healthy() {
			return nil
		}
		pool, err := pool.Connect()
		if err != nil {
			return err
		}
		g.conns.Store(addr, pool)
		return nil
	}
	pool, err = NewPool(ctx, addr, g.poolSize, append(g.dopts, dopts...)...)
	if err != nil {
		return err
	}
	atomic.AddUint64(&g.clientCount, 1)
	g.conns.Store(addr, pool)
	return nil
}

func (g *gRPCClient) Disconnect(addr string) error {
	pool, ok := g.conns.Load(addr)
	if !ok {
		return errors.ErrGRPCClientConnNotFound(addr)
	}
	g.conns.Delete(addr)
	atomic.AddUint64(&g.clientCount, ^uint64(0))
	if pool != nil {
		return pool.Disconnect()
	}
	return nil
}

func (g *gRPCClient) Close() error {
	closeList := make([]string, 0, int(atomic.LoadUint64(&g.clientCount)))
	g.conns.Range(func(addr string, pool *ClientConnPool) bool {
		if pool != nil {
			closeList = append(closeList, addr)
		}
		return true
	})
	for _, addr := range closeList {
		g.Disconnect(addr)
	}
	return nil
}

func (g *gRPCClient) GetAddrs() (connected []string, disconnected []string) {
	g.conns.Range(func(addr string, pool *ClientConnPool) bool {
		if pool.Healthy() {
			connected = append(connected, addr)
		} else {
			disconnected = append(disconnected, addr)
		}
		return true
	})
	return
}

func isHealthy(conn *ClientConn) bool {
	return conn != nil &&
		conn.GetState() != connectivity.Shutdown &&
		conn.GetState() != connectivity.TransientFailure
}
