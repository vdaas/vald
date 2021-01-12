//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/singleflight"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
)

type (
	Server       = grpc.Server
	ServerOption = grpc.ServerOption
	CallOption   = grpc.CallOption
	DialOption   = pool.DialOption
	ClientConn   = pool.ClientConn
)

type Client interface {
	StartConnectionMonitor(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, addr string, dopts ...DialOption) (pool.Conn, error)
	Disconnect(ctx context.Context, addr string) error
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
	RoundRobin(ctx context.Context, f func(ctx context.Context,
		conn *ClientConn,
		copts ...CallOption) (interface{}, error)) (interface{}, error)
	GetDialOption() []DialOption
	GetCallOption() []CallOption
	ConnectedAddrs() []string
	Close(ctx context.Context) error
}

type gRPCClient struct {
	addrs               map[string]struct{}
	atomicAddrs         AtomicAddrs
	poolSize            uint64
	clientCount         uint64
	conns               grpcConns
	hcDur               time.Duration
	prDur               time.Duration
	enablePoolRebalance bool
	resolveDNS          bool
	dopts               []DialOption
	copts               []CallOption
	roccd               string // reconnection old connection closing duration
	eg                  errgroup.Group
	bo                  backoff.Backoff
	gbo                 gbackoff.Config // grpc's original backoff configuration
	mcd                 time.Duration   // minimum connection timeout duration
	group               singleflight.Group
	crl                 sync.Map // connection request list
}

func New(opts ...Option) (c Client) {
	g := &gRPCClient{
		group: singleflight.New(),
		addrs: make(map[string]struct{}),
	}

	for _, opt := range append(defaultOptions, opts...) {
		opt(g)
	}
	g.atomicAddrs = newAddr(g.addrs)
	g.dopts = append(g.dopts, grpc.WithConnectParams(
		grpc.ConnectParams{
			Backoff: gbackoff.Config{
				MaxDelay:   g.gbo.MaxDelay,
				BaseDelay:  g.gbo.BaseDelay,
				Multiplier: g.gbo.Multiplier,
				Jitter:     g.gbo.Jitter,
			},
			MinConnectTimeout: g.mcd,
		},
	))
	return g
}

func (g *gRPCClient) StartConnectionMonitor(ctx context.Context) (<-chan error, error) {
	addrs, ok := g.atomicAddrs.GetAll()
	if !ok {
		return nil, errors.ErrGRPCTargetAddrNotFound
	}

	ech := make(chan error, len(addrs))
	for _, addr := range addrs {
		if len(addr) != 0 {
			_, err := g.Connect(ctx, addr, grpc.WithBlock())
			if err != nil {
				log.Error(err)
				ech <- err
			}
		}
	}

	if len(addrs) != 0 && atomic.LoadUint64(&g.clientCount) == 0 {
		return nil, errors.ErrGRPCClientConnNotFound(strings.Join(addrs, ",\t"))
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
		defer g.Close(context.Background())
		defer hcTick.Stop()
		defer prTick.Stop()
		for {
			select {
			case <-ctx.Done():
				if err != nil {
					return errors.Wrap(ctx.Err(), err.Error())
				}
				return ctx.Err()
			case <-prTick.C:
				if g.enablePoolRebalance {
					g.conns.Range(func(addr string, p pool.Conn) bool {
						if len(addr) != 0 && p != nil {
							var err error
							p, err = p.Reconnect(ctx, true)
							if err != nil {
								log.Error(err)
								ech <- err
								err = g.Disconnect(ctx, addr)
								if err != nil {
									log.Error(err)
									ech <- err
								}
							}
							if err == nil {
								g.conns.Store(addr, p)
							} else {
								g.conns.Delete(addr)
								g.atomicAddrs.Delete(addr)
							}
						}
						return true
					})
				}
			case <-hcTick.C:
				g.conns.Range(func(addr string, p pool.Conn) bool {
					if len(addr) != 0 && !p.IsHealthy(ctx) {
						var err error
						p, err = p.Reconnect(ctx, false)
						if err != nil {
							log.Error(err)
							ech <- err
							err = g.Disconnect(ctx, addr)
							if err != nil {
								log.Error(err)
								ech <- err
							}
						}
						if err == nil {
							g.conns.Store(addr, p)
						} else {
							g.conns.Delete(addr)
							g.atomicAddrs.Delete(addr)
						}
					}
					return true
				})
			}
			g.crl.Range(func(addr, _ interface{}) bool {
				defer g.crl.Delete(addr)
				_, err := g.Connect(ctx, addr.(string))
				if err != nil {
					log.Error(err)
				}
				return true
			})
		}
	}))
	return ech, nil
}

func (g *gRPCClient) Range(ctx context.Context,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.Range")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	g.conns.Range(func(addr string, p pool.Conn) bool {
		ssctx, sspan := trace.StartSpan(sctx, "vald/internal/grpc/Client.Range/"+addr)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		select {
		case <-ctx.Done():
			return false
		default:
			_, err := g.do(ssctx, p, addr, true, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption) (interface{}, error) {
				return nil, f(ictx, addr, conn, copts...)
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
	concurrency int, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) error {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.RangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	eg, egctx := errgroup.New(sctx)
	eg.Limitation(concurrency)
	g.conns.Range(func(addr string, p pool.Conn) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			ssctx, sspan := trace.StartSpan(sctx, "vald/internal/grpc/Client.RangeConcurrent/"+addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			select {
			case <-egctx.Done():
				return nil
			default:
				_, err := g.do(ssctx, p, addr, true, func(ictx context.Context,
					conn *ClientConn, copts ...CallOption) (interface{}, error) {
					return nil, f(ictx, addr, conn, copts...)
				})
				return err
			}
		}))
		return true
	})
	return eg.Wait()
}

func (g *gRPCClient) OrderedRange(ctx context.Context,
	orders []string, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.OrderedRange")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if orders == nil {
		log.Warn("no order found for OrderedRange")
		return g.Range(sctx, f)
	}
	for _, addr := range orders {
		select {
		case <-sctx.Done():
			return nil
		default:
			p, ok := g.conns.Load(addr)
			if !ok {
				g.crl.Store(addr, struct{}{})
				log.Warn(errors.ErrGRPCClientConnNotFound(addr))
				continue
			}
			ssctx, span := trace.StartSpan(sctx, "vald/internal/grpc/Client.OrderedRange/"+addr)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			_, err := g.do(ssctx, p, addr, true, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption) (interface{}, error) {
				return nil, f(ictx, addr, conn, copts...)
			})
			if err != nil {
				rerr = errors.Wrap(rerr, errors.ErrRPCCallFailed(addr, err).Error())
			}
		}
	}
	return rerr
}

func (g *gRPCClient) OrderedRangeConcurrent(ctx context.Context,
	orders []string, concurrency int, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error) (err error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.OrderedRangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if orders == nil {
		log.Warn("no order found for OrderedRangeConcurrent")
		return g.RangeConcurrent(sctx, concurrency, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.Limitation(concurrency)
	for _, order := range orders {
		addr := order
		eg.Go(safety.RecoverFunc(func() (err error) {
			p, ok := g.conns.Load(addr)
			if !ok {
				g.crl.Store(addr, struct{}{})
				return nil

			}
			ssctx, sspan := trace.StartSpan(sctx, "vald/internal/grpc/Client.OrderedRangeConcurrent/"+addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			select {
			case <-egctx.Done():
				return nil
			default:
				_, err := g.do(ssctx, p, addr, true, func(ictx context.Context,
					conn *ClientConn, copts ...CallOption) (interface{}, error) {
					return nil, f(ictx, addr, conn, copts...)
				})
				return err
			}
		}))
	}
	return eg.Wait()
}

func (g *gRPCClient) RoundRobin(ctx context.Context, f func(ctx context.Context,
	conn *ClientConn, copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.RoundRobin")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.bo != nil {
		return g.bo.Do(sctx, func(ictx context.Context) (r interface{}, ret bool, err error) {
			addr, ok := g.atomicAddrs.Next()
			if !ok {
				return nil, false, errors.ErrGRPCClientNotFound
			}
			p, ok := g.conns.Load(addr)
			if !ok {
				g.crl.Store(addr, struct{}{})
				return nil, true, errors.ErrGRPCClientConnNotFound(addr)
			}
			r, err = g.do(ictx, p, addr, false, f)
			if err != nil {
				return nil, true, err
			}
			return r, false, nil
		})
	}
	addr, ok := g.atomicAddrs.Next()
	if !ok {
		return nil, errors.ErrGRPCClientNotFound
	}
	return g.Do(sctx, addr, f)
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.Do/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	p, ok := g.conns.Load(addr)
	if !ok {
		g.crl.Store(addr, struct{}{})
		return nil, errors.ErrGRPCClientConnNotFound(addr)
	}
	return g.do(sctx, p, addr, true, f)
}

func (g *gRPCClient) do(ctx context.Context, p pool.Conn, addr string, enableBackoff bool,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	sctx, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.do/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.bo != nil && enableBackoff {
		data, err = g.bo.Do(sctx, func(ictx context.Context) (r interface{}, ret bool, err error) {
			err = p.Do(func(conn *ClientConn) (err error) {
				if conn == nil {
					return errors.ErrGRPCClientConnNotFound(addr)
				}
				r, err = f(ictx, conn, g.copts...)
				return err
			})
			if err != nil {
				return nil, err != nil, err
			}
			return r, false, nil
		})
	} else {
		err = p.Do(func(conn *ClientConn) (err error) {
			if conn == nil {
				return errors.ErrGRPCClientConnNotFound(addr)
			}
			data, err = f(sctx, conn, g.copts...)
			return err
		})
	}
	if err != nil {
		if p.Len() <= 0 {
			g.conns.Delete(addr)
		}
		return nil, errors.ErrRPCCallFailed(addr, err)
	}
	return data, nil
}

func (g *gRPCClient) GetDialOption() []DialOption {
	return g.dopts
}

func (g *gRPCClient) GetCallOption() []CallOption {
	return g.copts
}

func (g *gRPCClient) Connect(ctx context.Context, addr string, dopts ...DialOption) (conn pool.Conn, err error) {
	_, span := trace.StartSpan(ctx, "vald/internal/grpc/Client.Connect/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, _, err = g.group.Do(ctx, "connect-"+addr, func() (interface{}, error) {
		conn, ok := g.conns.Load(addr)
		if ok && conn != nil {
			if conn.IsHealthy(ctx) {
				g.atomicAddrs.Add(addr)
				return conn, nil
			}
			log.Debugf("connecting unhealthy pool addr= %s", addr)
			conn, err = conn.Connect(ctx)
			if err == nil {
				g.conns.Store(addr, conn)
				g.atomicAddrs.Add(addr)
				return conn, nil
			}
			log.Warnf("failed to reconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
			g.conns.Delete(addr)
			atomic.AddUint64(&g.clientCount, ^uint64(0))
			if conn != nil {
				err = conn.Disconnect()
				if err != nil {
					log.Warnf("failed to disconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
					g.conns.Delete(addr)
				}
			}
		} else if conn == nil {
			g.conns.Delete(addr)
		} else {
			err = conn.Disconnect()
			if err != nil {
				log.Warnf("failed to disconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
			}
			g.conns.Delete(addr)
		}

		log.Warnf("creating new connection pool for addr = %s", addr)
		opts := []pool.Option{
			pool.WithAddr(addr),
			pool.WithSize(g.poolSize),
			pool.WithDialOptions(append(g.dopts, dopts...)...),
			pool.WithResolveDNS(g.resolveDNS),
		}
		if g.bo != nil {
			opts = append(opts, pool.WithBackoff(g.bo))
		}
		conn, err = pool.New(ctx, opts...)
		if err != nil {
			g.conns.Delete(addr)
			return nil, err
		}
		log.Warnf("connecting to new connection pool for addr= %s", addr)
		conn, err = conn.Connect(ctx)
		if err != nil {
			g.conns.Delete(addr)
			return nil, err
		}
		atomic.AddUint64(&g.clientCount, 1)
		g.conns.Store(addr, conn)
		g.atomicAddrs.Add(addr)
		return conn, nil
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (g *gRPCClient) Disconnect(ctx context.Context, addr string) error {
	_, _, err := g.group.Do(ctx, "disconnect-"+addr, func() (interface{}, error) {
		p, ok := g.conns.Load(addr)
		if !ok {
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
		g.conns.Delete(addr)
		g.atomicAddrs.Delete(addr)
		atomic.AddUint64(&g.clientCount, ^uint64(0))
		if p != nil {
			log.Debugf("disconnecting grpc client conn pool addr= %s", addr)
			return nil, p.Disconnect()
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *gRPCClient) ConnectedAddrs() []string {
	addrs, _ := g.atomicAddrs.GetAll()
	return addrs
}

func (g *gRPCClient) Close(ctx context.Context) (err error) {
	var closeList []string
	if cc := int(atomic.LoadUint64(&g.clientCount)); cc > 0 {
		closeList = make([]string, 0, cc)
	}
	g.conns.Range(func(addr string, pool pool.Conn) bool {
		if pool != nil {
			closeList = append(closeList, addr)
		}
		return true
	})
	for _, addr := range closeList {
		derr := g.Disconnect(ctx, addr)
		if derr != nil {
			err = errors.Wrap(err, derr.Error())
		}
	}
	return err
}
