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
type ServerOption = grpc.ServerOption
type ClientConn = grpc.ClientConn
type CallOption = grpc.CallOption
type DialOption = grpc.DialOption

type Client interface {
	StartConnectionMonitor(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, addr string, dopts ...DialOption) error
	Disconnect(addr string) error
	Range(ctx context.Context,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	RangeConcurrent(ctx context.Context,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	OrderedRange(ctx context.Context,
		order []string,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	OrderedRangeConcurrent(ctx context.Context,
		order []string,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	Do(ctx context.Context, addr string,
		f func(ctx context.Context,
			conn *ClientConn,
			copts ...CallOption) (interface{}, error)) (interface{}, error)
	GetDialOption() []DialOption
	GetCallOption() []CallOption
	Close() error
}

type gRPCClient struct {
	addrs               []string
	poolSize            uint64
	clientCount         uint64
	conns               grpcConns
	hcDur               time.Duration
	prDur               time.Duration
	enablePoolRebalance bool
	dopts               []DialOption
	copts               []CallOption
	eg                  errgroup.Group
	bo                  backoff.Backoff
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
		prTick := &time.Ticker{
			C: make(chan time.Time),
		}
		if g.enablePoolRebalance {
			prTick.Stop()
			prTick = time.NewTicker(g.prDur)
		}
		hcTick := time.NewTicker(g.hcDur)
		defer close(ech)
		defer g.Close()
		defer hcTick.Stop()
		defer prTick.Stop()
		for {
			select {
			case <-ctx.Done():
				err = g.Close()
				if err != nil {
					return errors.Wrap(ctx.Err(), err.Error())
				}
				return ctx.Err()
			case <-prTick.C:
				if g.enablePoolRebalance {
					g.conns.Range(func(addr string, pool ClientConnPool) bool {
						if len(addr) != 0 && !pool.IsHealthy() {
							err = g.Disconnect(addr)
							if err != nil {
								log.Error(err)
							}
						} else {
							cp, err := pool.Reconnect(ctx, false)
							if err != nil {
								log.Error(err)
							} else {
								g.conns.Store(addr, cp)
							}
						}
						return true
					})
				}
			case <-hcTick.C:
				reconnList := make([]string, 0, int(atomic.LoadUint64(&g.clientCount)))
				g.conns.Range(func(addr string, pool ClientConnPool) bool {
					if len(addr) != 0 && !pool.IsHealthy() {
						reconnList = append(reconnList, addr)
					}
					return true
				})
				for _, addr := range reconnList {
					err = g.Disconnect(addr)
					if err != nil {
						log.Error(err)
						ech <- err
					}
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
						err = g.Disconnect(addr)
						if err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	}))
	return ech, nil
}

func (g *gRPCClient) Range(ctx context.Context,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	g.conns.Range(func(addr string, pool ClientConnPool) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			var err error
			if g.bo != nil {
				_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
					return nil, pool.Do(func(conn *ClientConn) (err error) {
						if conn == nil {
							return errors.ErrGRPCClientConnNotFound(addr)
						}
						return f(ctx, addr, conn, g.copts...)
					})
				})
			} else {
				err = pool.Do(func(conn *ClientConn) (err error) {
					if conn == nil {
						return errors.ErrGRPCClientConnNotFound(addr)
					}
					return f(ctx, addr, conn, g.copts...)
				})
			}
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
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) error {
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	g.conns.Range(func(addr string, pool ClientConnPool) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			select {
			case <-egctx.Done():
				return nil
			default:
				if g.bo != nil {
					_, err = g.bo.Do(egctx, func() (r interface{}, err error) {
						return nil, pool.Do(func(conn *ClientConn) (err error) {
							if conn == nil {
								return errors.ErrGRPCClientConnNotFound(addr)
							}
							return f(egctx, addr, conn, g.copts...)
						})
					})
				} else {
					err = pool.Do(func(conn *ClientConn) (err error) {
						if conn == nil {
							return errors.ErrGRPCClientConnNotFound(addr)
						}
						return f(egctx, addr, conn, g.copts...)
					})
				}
				if err != nil {
					return errors.ErrRPCCallFailed(addr, err)
				}
				return nil
			}
		}))
		return true
	})
	return eg.Wait()
}

func (g *gRPCClient) OrderedRange(ctx context.Context,
	orders []string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	if orders == nil {
		log.Warn("no order found for OrderedRange")
		return g.Range(ctx, f)
	}
	var err error
	for _, addr := range orders {
		select {
		case <-ctx.Done():
			return nil
		default:
			pool, ok := g.conns.Load(addr)
			if ok {
				if g.bo != nil {
					_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
						return nil, pool.Do(func(conn *ClientConn) (err error) {
							if conn == nil {
								return errors.ErrGRPCClientConnNotFound(addr)
							}
							return f(ctx, addr, conn, g.copts...)
						})
					})
				} else {
					err = pool.Do(func(conn *ClientConn) (err error) {
						if conn == nil {
							return errors.ErrGRPCClientConnNotFound(addr)
						}
						return f(ctx, addr, conn, g.copts...)
					})
				}
				if err != nil {
					rerr = errors.Wrap(rerr, errors.ErrRPCCallFailed(addr, err).Error())
				}
			} else {
				log.Warnf("connection %s not found for OrderedRange", addr)
			}
		}
	}
	return rerr
}

func (g *gRPCClient) OrderedRangeConcurrent(ctx context.Context,
	orders []string,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (err error) {
	if orders == nil {
		log.Warn("no order found for OrderedRangeConcurrent")
		return g.RangeConcurrent(ctx, concurrency, f)
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	for _, order := range orders {
		addr := order
		pool, ok := g.conns.Load(addr)
		if ok {
			eg.Go(safety.RecoverFunc(func() (err error) {
				select {
				case <-egctx.Done():
					return nil
				default:
					if g.bo != nil {
						_, err = g.bo.Do(egctx, func() (r interface{}, err error) {
							return nil, pool.Do(func(conn *ClientConn) (err error) {
								if conn == nil {
									return errors.ErrGRPCClientConnNotFound(addr)
								}
								return f(egctx, addr, conn, g.copts...)
							})
						})
					} else {
						err = pool.Do(func(conn *ClientConn) (err error) {
							if conn == nil {
								return errors.ErrGRPCClientConnNotFound(addr)
							}
							return f(egctx, addr, conn, g.copts...)
						})
					}
					if err != nil {
						return errors.ErrRPCCallFailed(addr, err)
					}
					return nil
				}
			}))
		} else {
			log.Warnf("connection %s not found for OrderedRangeConcurrent", addr)
		}
	}
	return eg.Wait()
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	pool, ok := g.conns.Load(addr)
	if !ok {
		return nil, errors.ErrGRPCClientConnNotFound(addr)
	}
	if g.bo != nil {
		data, err = g.bo.Do(ctx, func() (r interface{}, err error) {
			err = pool.Do(func(conn *ClientConn) (err error) {
				if conn == nil {
					return errors.ErrGRPCClientConnNotFound(addr)
				}
				r, err = f(ctx, conn, g.copts...)
				return err
			})
			if err != nil {
				return nil, err
			}
			return r, err
		})
	} else {
		err = pool.Do(func(conn *ClientConn) (err error) {
			if conn == nil {
				return errors.ErrGRPCClientConnNotFound(addr)
			}
			data, err = f(ctx, conn, g.copts...)
			return err
		})
	}
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
	if ok && pool != nil {
		if pool.IsHealthy() {
			return nil
		}
		log.Debugf("connecting unhealthy pool addr= %s", addr)
		pool, err = pool.Connect(ctx)
		if err == nil {
			g.conns.Store(addr, pool)
			return nil
		}
		log.Warnf("failed to reconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
		g.conns.Delete(addr)
		atomic.AddUint64(&g.clientCount, ^uint64(0))
		err = pool.Disconnect()
		if err != nil {
			log.Warnf("failed to disconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
			g.conns.Delete(addr)
		}
	} else if pool == nil {
		g.conns.Delete(addr)
	}

	log.Warnf("creating new connection pool for addr= %s", addr)
	pool, err = NewPool(ctx, addr, g.poolSize, append(g.dopts, dopts...)...)
	if err != nil {
		g.conns.Delete(addr)
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
	g.conns.Range(func(addr string, pool ClientConnPool) bool {
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

func isHealthy(conn *ClientConn) bool {
	return conn != nil &&
		conn.GetState() != connectivity.Shutdown &&
		conn.GetState() != connectivity.TransientFailure
}
