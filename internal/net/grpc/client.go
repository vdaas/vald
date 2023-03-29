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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/singleflight"
	"github.com/vdaas/vald/internal/strings"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
)

type (
	CallOption = grpc.CallOption
	DialOption = pool.DialOption
	ClientConn = pool.ClientConn
)

type Client interface {
	StartConnectionMonitor(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, addr string, dopts ...DialOption) (pool.Conn, error)
	IsConnected(ctx context.Context, addr string) bool
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
	GetBackoff() backoff.Backoff
	ConnectedAddrs() []string
	Close(ctx context.Context) error
}

type gRPCClient struct {
	addrs               map[string]struct{}
	poolSize            uint64
	clientCount         uint64
	conns               grpcConns
	hcDur               time.Duration
	prDur               time.Duration
	dialer              net.Dialer
	enablePoolRebalance bool
	resolveDNS          bool
	dopts               []DialOption
	copts               []CallOption
	roccd               string // reconnection old connection closing duration
	eg                  errgroup.Group
	bo                  backoff.Backoff
	cb                  circuitbreaker.CircuitBreaker
	gbo                 gbackoff.Config // grpc's original backoff configuration
	mcd                 time.Duration   // minimum connection timeout duration
	group               singleflight.Group
	crl                 sync.Map // connection request list

	ech            <-chan error
	monitorRunning atomic.Bool
	stopMonitor    context.CancelFunc
}

const apiName = "vald/internal/net/grpc"

func New(opts ...Option) (c Client) {
	g := &gRPCClient{
		group: singleflight.New(),
		addrs: make(map[string]struct{}),
	}

	for _, opt := range append(defaultOptions, opts...) {
		opt(g)
	}
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
	g.monitorRunning.Store(false)
	return g
}

func (g *gRPCClient) StartConnectionMonitor(ctx context.Context) (<-chan error, error) {
	logger.Init()
	if g.monitorRunning.Load() {
		return g.ech, nil
	}
	g.monitorRunning.Store(true)

	addrs := make([]string, len(g.addrs))
	for addr := range g.addrs {
		addrs = append(addrs, addr)
	}

	if g.dialer != nil {
		g.dialer.StartDialerCache(ctx)
	}

	ech := make(chan error, len(addrs))
	for _, addr := range addrs {
		if len(addr) != 0 {
			_, err := g.Connect(ctx, addr, grpc.WithBlock())
			if err != nil {
				if !errors.Is(err, context.Canceled) &&
					!errors.Is(err, context.DeadlineExceeded) &&
					!errors.Is(err, errors.ErrCircuitBreakerOpenState) &&
					!errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) &&
					!errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) &&
					!errors.Is(err, errors.ErrGRPCClientNotFound) {
					log.Errorf("failed to initial connection to %s,\terror: %v", addr, err)
					ech <- err
				} else {
					log.Warn(err)
				}
			}
		}
	}

	if len(addrs) != 0 && atomic.LoadUint64(&g.clientCount) == 0 {
		err := errors.ErrGRPCClientConnNotFound(strings.Join(addrs, ",\t"))
		log.Error(err)
		return nil, err
	}

	ctx, g.stopMonitor = context.WithCancel(ctx)
	g.eg.Go(safety.RecoverFunc(func() (err error) {
		defer g.monitorRunning.Store(false)
		prTick := &time.Ticker{ // pool rebalance ticker
			C: make(chan time.Time),
		}
		// this duration is for timeout to prevent blocking health check loop, which should be minimum duration of hcDur and prDur
		reconnLimitDuration := time.Second

		hcTick := time.NewTicker(g.hcDur) // health check ticker
		if g.enablePoolRebalance && g.prDur.Nanoseconds() > 0 {
			prTick.Stop()
			prTick = time.NewTicker(g.prDur)
			reconnLimitDuration = time.Duration(int64(math.Min(float64(g.hcDur.Nanoseconds()), float64(g.prDur.Nanoseconds()))))
		} else {
			reconnLimitDuration = g.hcDur
		}
		defer close(ech)
		defer g.Close(context.Background())
		defer hcTick.Stop()
		defer prTick.Stop()
		disconnectTargets := make([]string, 0, len(addrs))
		for {
			select {
			case <-ctx.Done():
				if err != nil {
					return errors.Wrap(ctx.Err(), err.Error())
				}
				return ctx.Err()
			case <-prTick.C:
				if g.enablePoolRebalance {
					err = g.conns.Range(func(addr string, p pool.Conn) bool {
						// if addr or pool is nil or empty the registration of conns is invalid let's disconnect them
						if len(addr) == 0 || p == nil {
							disconnectTargets = append(disconnectTargets, addr)
							return true
						}
						var err error
						// for rebalancing connection we don't need to check connection health
						p, err = p.Connect(ctx)
						if err != nil {
							if !errors.Is(err, context.Canceled) &&
								!errors.Is(err, context.DeadlineExceeded) &&
								!errors.Is(err, errors.ErrCircuitBreakerOpenState) &&
								!errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) &&
								!errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) &&
								!errors.Is(err, errors.ErrGRPCClientNotFound) {
								log.Error(err)
								ech <- err
							} else {
								log.Warn(err)
							}
						}
						// if rebalanced connection pool is nil even error is nil we should disconnect and delete it
						if err == nil && p == nil {
							disconnectTargets = append(disconnectTargets, addr)
							return true
						}
						// if connection pool could not recover we should try next connection loop
						if err != nil || !p.IsHealthy(ctx) {
							g.crl.Store(addr, false)
							return true
						}
						g.conns.Store(addr, p)
						return true
					})
				}
			case <-hcTick.C:
				err = g.conns.Range(func(addr string, p pool.Conn) bool {
					// if addr or pool is nil or empty the registration of conns is invalid let's disconnect them
					if len(addr) == 0 || p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					// for health check we don't need to reconnect when pool is healthy
					if p.IsHealthy(ctx) {
						return true
					}
					var err error
					// if not healthy we should try reconnect
					p, err = p.Reconnect(ctx, false)
					if err != nil {
						if !errors.Is(err, context.Canceled) &&
							!errors.Is(err, context.DeadlineExceeded) &&
							!errors.Is(err, errors.ErrCircuitBreakerOpenState) &&
							!errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) &&
							!errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) &&
							!errors.Is(err, errors.ErrGRPCClientNotFound) {
							log.Error(err)
							ech <- err
						} else {
							log.Warn(err)
						}
					}
					// if rebalanced connection pool is nil even error is nil we should disconnect and delete it
					if err == nil && p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					// if connection pool could not recover we should try next connection loop
					if err != nil || !p.IsHealthy(ctx) {
						g.crl.Store(addr, false)
						return true
					}
					g.conns.Store(addr, p)
					return true
				})
			}
			if err != nil && errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) && len(addrs) != 0 {
				for _, addr := range addrs {
					if len(addr) != 0 {
						log.Debugf("connection for %s not found in connection map will re-connect soon", addr)
						g.crl.Store(addr, false)
					}
				}
			}
			clctx, cancel := context.WithTimeout(ctx, reconnLimitDuration)
			g.crl.Range(func(a, bo interface{}) bool {
				select {
				case <-clctx.Done():
					return false
				default:
					defer g.crl.Delete(a)
					addr, ok := a.(string)
					if !ok {
						return true
					}
					var p pool.Conn
					if enabled, ok := bo.(bool); ok && enabled && g.bo != nil {
						_, err = g.bo.Do(clctx, func(ictx context.Context) (r interface{}, ret bool, err error) {
							p, err = g.Connect(ictx, addr)
							return nil, err != nil, err
						})
					} else {
						p, err = g.Connect(clctx, addr)
					}
					if err != nil || p == nil || !p.IsHealthy(ctx) {
						log.Debugf("connection for %s is not healthy will delete soon,\terror: %v,\tpool: [%v]", addr, err, p)
						disconnectTargets = append(disconnectTargets, addr)
					} else {
						g.conns.Store(addr, p)
					}
					return true
				}
			})
			cancel()
			disconnected := make(map[string]bool, len(disconnectTargets))
			for _, addr := range disconnectTargets {
				if !disconnected[addr] {
					err = g.Disconnect(ctx, addr)
					if err != nil {
						if !errors.Is(err, context.Canceled) &&
							!errors.Is(err, context.DeadlineExceeded) &&
							!errors.Is(err, errors.ErrCircuitBreakerOpenState) &&
							!errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) &&
							!errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) &&
							!errors.Is(err, errors.ErrGRPCClientNotFound) {
							log.Error(err)
							ech <- err
						} else {
							log.Warn(err)
						}
					}
					disconnected[addr] = true
				}
			}
			disconnectTargets = disconnectTargets[:0]
		}
	}))
	g.ech = ech
	return ech, nil
}

func (g *gRPCClient) Range(ctx context.Context,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (rerr error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.Range")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	return g.conns.Range(func(addr string, p pool.Conn) bool {
		ssctx, sspan := trace.StartSpan(sctx, apiName+"/Client.Range/"+addr)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		select {
		case <-ctx.Done():
			return false
		default:
			g.connectWithBackoff(ssctx, p, addr, true, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption,
			) (interface{}, error) {
				return nil, f(ictx, addr, conn, copts...)
			})
		}
		return true
	})
}

func (g *gRPCClient) RangeConcurrent(ctx context.Context,
	concurrency int, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (rerr error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.RangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if concurrency == 0 || concurrency == 1 {
		return g.Range(ctx, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.Limitation(concurrency)
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	rerr = g.conns.Range(func(addr string, p pool.Conn) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			ssctx, sspan := trace.StartSpan(sctx, apiName+"/Client.RangeConcurrent/"+addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			select {
			case <-egctx.Done():
				return nil
			default:
				g.connectWithBackoff(ssctx, p, addr, true, func(ictx context.Context,
					conn *ClientConn, copts ...CallOption,
				) (interface{}, error) {
					err := f(ictx, addr, conn, copts...)
					return nil, err
				})
				return nil
			}
		}))
		return true
	})
	if rerr != nil {
		return rerr
	}
	rerr = eg.Wait()
	return rerr
}

func (g *gRPCClient) OrderedRange(ctx context.Context,
	orders []string, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (rerr error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.OrderedRange")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRange")
		return g.Range(sctx, f)
	}
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	for _, addr := range orders {
		select {
		case <-sctx.Done():
			return nil
		default:
			p, ok := g.conns.Load(addr)
			if !ok || p == nil {
				g.crl.Store(addr, true)
				log.Warnf("gRPCClient.OrderedRange operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, errors.ErrGRPCClientConnNotFound(addr))
				continue
			}
			ssctx, span := trace.StartSpan(sctx, apiName+"/Client.OrderedRange/"+addr)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			g.connectWithBackoff(ssctx, p, addr, true, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption,
			) (interface{}, error) {
				return nil, f(ictx, addr, conn, copts...)
			})
		}
	}
	return nil
}

func (g *gRPCClient) OrderedRangeConcurrent(ctx context.Context,
	orders []string, concurrency int, f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (rerr error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.OrderedRangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRangeConcurrent")
		return g.RangeConcurrent(sctx, concurrency, f)
	}
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if concurrency < 2 {
		return g.OrderedRange(ctx, orders, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.Limitation(concurrency)
	for _, order := range orders {
		addr := order
		eg.Go(safety.RecoverFunc(func() (err error) {
			p, ok := g.conns.Load(addr)
			if !ok || p == nil {
				g.crl.Store(addr, true)
				log.Warnf("gRPCClient.OrderedRangeConcurrent operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, errors.ErrGRPCClientConnNotFound(addr))
				return nil
			}
			ssctx, sspan := trace.StartSpan(sctx, apiName+"/Client.OrderedRangeConcurrent/"+addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			select {
			case <-egctx.Done():
				return nil
			default:
				g.connectWithBackoff(ssctx, p, addr, true, func(ictx context.Context,
					conn *ClientConn, copts ...CallOption,
				) (interface{}, error) {
					return nil, f(ictx, addr, conn, copts...)
				})
				return nil
			}
		}))
	}
	return eg.Wait()
}

func (g *gRPCClient) RoundRobin(ctx context.Context, f func(ctx context.Context,
	conn *ClientConn, copts ...CallOption) (interface{}, error),
) (data interface{}, err error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.RoundRobin")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.conns.Len() == 0 {
		return nil, errors.ErrGRPCClientConnNotFound("*")
	}

	var boName string
	if boName = FromGRPCMethod(sctx); len(boName) != 0 {
		sctx = backoff.WithBackoffName(sctx, boName)
	}

	do := func() (data interface{}, err error) {
		cerr := g.conns.Range(func(addr string, p pool.Conn) bool {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return false
			default:
				if p != nil && p.IsHealthy(sctx) {
					ctx, span := trace.StartSpan(sctx, apiName+"/Client.RoundRobin/"+addr)
					defer func() {
						if span != nil {
							span.End()
						}
					}()
					var boName string
					ctx = WrapGRPCMethod(ctx, addr)
					if boName = FromGRPCMethod(ctx); len(boName) != 0 {
						ctx = backoff.WithBackoffName(ctx, boName)
					}
					if g.cb != nil && len(boName) > 0 {
						data, err = g.cb.Do(ctx, boName, func(ictx context.Context) (interface{}, error) {
							return g.connectWithBackoff(ictx, p, addr, false, f)
						})
						if err != nil {
							return true
						}
						return false
					}
					data, err = g.connectWithBackoff(ctx, p, addr, false, f)
					if err != nil {
						return true
					}
					return false
				}
				g.crl.Store(addr, true)
				return true
			}
		})
		if cerr != nil {
			return nil, cerr
		}
		return data, err
	}

	if g.bo != nil {
		return g.bo.Do(sctx, func(ictx context.Context) (r interface{}, ret bool, err error) {
			r, err = do()
			if err != nil {
				if errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, errors.ErrCircuitBreakerOpenState) ||
					errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) ||
					errors.Is(err, errors.ErrGRPCClientNotFound) {
					return nil, false, err
				}
				st, ok := status.FromError(err)
				if !ok || st == nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) {
						return nil, false, err
					}
					return nil, err != nil, err
				}
				status.Log(st.Code(), err)
				switch st.Code() {
				case codes.Internal,
					codes.Unavailable,
					codes.ResourceExhausted:
					return nil, err != nil, err
				}
				return nil, false, err
			}
			return r, false, nil
		})
	}
	return do()
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (interface{}, error),
) (data interface{}, err error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.Do/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	p, ok := g.conns.Load(addr)
	if !ok || p == nil {
		g.crl.Store(addr, true)
		err = errors.ErrGRPCClientConnNotFound(addr)
		log.Warnf("gRPCClient.Do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
		return nil, err
	}
	return g.connectWithBackoff(sctx, p, addr, true, f)
}

func (g *gRPCClient) connectWithBackoff(ctx context.Context, p pool.Conn, addr string, enableBackoff bool,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (interface{}, error),
) (data interface{}, err error) {
	if p == nil {
		g.crl.Store(addr, true)
		err = errors.ErrGRPCClientConnNotFound(addr)
		log.Warnf("gRPCClient.do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
		return nil, err
	}
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.do/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.bo != nil && enableBackoff {
		var boName string
		sctx = WrapGRPCMethod(sctx, addr)
		if boName = FromGRPCMethod(sctx); len(boName) != 0 {
			sctx = backoff.WithBackoffName(sctx, boName)
		}
		do := func(ctx context.Context) (r interface{}, ret bool, err error) {
			err = p.Do(ctx, func(conn *ClientConn) (err error) {
				if conn == nil {
					return errors.ErrGRPCClientConnNotFound(addr)
				}
				r, err = f(ctx, conn, g.copts...)
				return err
			})
			if err != nil {
				if errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, errors.ErrCircuitBreakerOpenState) ||
					errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) ||
					errors.Is(err, errors.ErrGRPCClientNotFound) ||
					p.IsIPConn() && errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) {
					return nil, false, err
				}
				st, ok := status.FromError(err)
				if !ok || st == nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) {
						return nil, false, err
					}
					return nil, err != nil, err
				}
				status.Log(st.Code(), err)
				switch st.Code() {
				case codes.Internal,
					codes.Unavailable,
					codes.ResourceExhausted:
					return nil, err != nil, err
				}
				return nil, false, err
			}
			return r, false, nil
		}
		data, err = g.bo.Do(sctx, func(ictx context.Context) (r interface{}, ret bool, err error) {
			if g.cb != nil && len(boName) > 0 {
				r, err = g.cb.Do(ictx, boName, func(ictx context.Context) (interface{}, error) {
					r, ret, err = do(ictx)
					if err != nil && !ret {
						return r, errors.NewErrCircuitBreakerIgnorable(err)
					}
					return r, err
				})
				if err != nil {
					if errors.Is(err, context.Canceled) ||
						errors.Is(err, context.DeadlineExceeded) ||
						errors.Is(err, errors.ErrCircuitBreakerOpenState) ||
						errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) ||
						errors.Is(err, errors.ErrGRPCClientNotFound) {
						return nil, false, err
					}
				}
				return r, ret, err
			}
			return do(ictx)
		})
	} else {
		err = p.Do(sctx, func(conn *ClientConn) (err error) {
			if conn == nil {
				return errors.ErrGRPCClientConnNotFound(addr)
			}
			data, err = f(sctx, conn, g.copts...)
			return err
		})
	}
	if err != nil {
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

func (g *gRPCClient) GetBackoff() backoff.Backoff {
	return g.bo
}

func (g *gRPCClient) Connect(ctx context.Context, addr string, dopts ...DialOption) (conn pool.Conn, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Connect/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ci, shared, err := g.group.Do(ctx, "connect-"+addr, func() (interface{}, error) {
		var ok bool
		conn, ok = g.conns.Load(addr)
		if ok && conn != nil {
			if conn.IsHealthy(ctx) {
				return conn, nil
			}
			log.Debugf("connecting unhealthy pool addr= %s", addr)
			conn, err = conn.Connect(ctx)
			if err == nil && conn != nil && conn.IsHealthy(ctx) {
				g.conns.Store(addr, conn)
				return conn, nil
			}
			log.Warnf("failed to reconnect unhealthy pool conn=[%v]\terror= %v\t trying to disconnect", conn, err)
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
		if err != nil || conn == nil {
			derr := g.Disconnect(ctx, addr)
			if derr != nil && !errors.Is(derr, errors.ErrGRPCClientConnNotFound(addr)) {
				log.Warnf("failed to disconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
				err = errors.Wrap(err, derr.Error())
			}
			return nil, err
		}
		log.Warnf("connecting to new connection pool for addr= %s", addr)
		conn, err = conn.Connect(ctx)
		if err != nil {
			log.Error(err)
			derr := g.Disconnect(ctx, addr)
			if derr != nil && !errors.Is(derr, errors.ErrGRPCClientConnNotFound(addr)) {
				log.Warnf("failed to disconnect unhealthy pool addr= %s\terror= %s", addr, err.Error())
				err = errors.Wrap(err, derr.Error())
			}
			return nil, err
		}
		if conn == nil || !conn.IsHealthy(ctx) {
			log.Debugf("connection to %s is nil or unhealthy, conn=%v", addr, conn)
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
		atomic.AddUint64(&g.clientCount, 1)
		g.conns.Store(addr, conn)
		return conn, nil
	})
	if err != nil {
		return nil, err
	}
	if shared {
		sconn, ok := ci.(pool.Conn)
		if ok {
			return sconn, nil
		}
	}
	return conn, nil
}

func (g *gRPCClient) IsConnected(ctx context.Context, addr string) bool {
	p, ok := g.conns.Load(addr)
	if !ok || p == nil {
		return false
	}
	return p.IsHealthy(ctx)
}

func (g *gRPCClient) Disconnect(ctx context.Context, addr string) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Disconnect/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, _, err := g.group.Do(ctx, "disconnect-"+addr, func() (interface{}, error) {
		p, ok := g.conns.Load(addr)
		if !ok || p == nil {
			g.conns.Delete(addr)
			log.Debugf("target %s's gRPC connection pool is already unavailable", addr)
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
		g.conns.Delete(addr)
		atomic.AddUint64(&g.clientCount, ^uint64(0))
		if p != nil {
			log.Debugf("gRPC client connection pool addr = %s will disconnect soon...", addr)
			return nil, p.Disconnect()
		}
		return nil, nil
	})
	if err != nil {
		g.conns.Delete(addr)
		return err
	}
	return nil
}

func (g *gRPCClient) ConnectedAddrs() (addrs []string) {
	addrs = make([]string, 0, g.conns.Len())
	err := g.conns.Range(func(addr string, p pool.Conn) bool {
		if p != nil && p.IsHealthy(context.Background()) {
			addrs = append(addrs, addr)
		}
		return true
	})
	if err != nil {
		return nil
	}
	return addrs
}

func (g *gRPCClient) Close(ctx context.Context) (err error) {
	if g.stopMonitor != nil {
		g.stopMonitor()
	}
	g.conns.Range(func(addr string, p pool.Conn) bool {
		derr := g.Disconnect(ctx, addr)
		if derr != nil && !errors.Is(derr, errors.ErrGRPCClientConnNotFound(addr)) {
			err = errors.Wrap(err, derr.Error())
		}
		return true
	})
	return err
}
