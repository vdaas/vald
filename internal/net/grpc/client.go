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
	clientCount uint64
	conns       gRPCConns
	hcDur       time.Duration
	dopts       []DialOption
	copts       []CallOption
	eg          errgroup.Group
	bo          backoff.Backoff
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
				g.conns.Range(func(addr string, conn *ClientConn) bool {
					if len(addr) != 0 && !g.isHealthy(conn) {
						reconnList = append(reconnList, addr)
					}
					return true
				})

				for _, addr := range reconnList {
					if g.bo != nil {
						_, err = g.bo.Do(ctx, func() (interface{}, error) {
							_, err := g.reconnect(ctx, addr, nil)
							return nil, err
						})
					} else {
						_, err = g.reconnect(ctx, addr, nil)
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
	g.conns.Range(func(addr string, conn *ClientConn) bool {
		wrapf := func(addr string, conn *ClientConn, copts ...CallOption) (err error) {
			conn, err = g.reconnect(ctx, addr, conn)
			if err != nil {
				return errors.Wrap(err, errors.ErrGRPCClientConnNotFound(addr).Error())
			}
			return f(addr, conn, copts...)
		}
		select {
		case <-ctx.Done():
			return false
		default:
			var err error
			if g.bo != nil {
				_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
					err = wrapf(addr, conn, g.copts...)
					return
				})
			} else {
				err = wrapf(addr, conn, g.copts...)
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
	f func(addr string, conn *ClientConn, copts ...CallOption) error) (rerr error) {
	wrapf := func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) (err error) {
		conn, err = g.reconnect(ctx, addr, conn)
		if err != nil {
			return errors.Wrap(err, errors.ErrGRPCClientConnNotFound(addr).Error())
		}
		return f(addr, conn, copts...)
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	g.conns.Range(func(addr string, conn *ClientConn) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			select {
			case <-egctx.Done():
				return nil
			default:
				var err error
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
			}
			return nil
		}))
		return true
	})
	return eg.Wait()
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(conn *ClientConn,
		copts ...CallOption) (interface{}, error)) (data interface{}, err error) {
	var conn *ClientConn
	wrapf := func(_ *ClientConn, copts ...CallOption) (ret interface{}, err error) {
		conn, err = g.reconnect(ctx, addr, conn)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrGRPCClientConnNotFound(addr).Error())
		}
		return f(conn, copts...)
	}
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
	conn, ok := g.conns.Load(addr)
	if ok {
		if g.isHealthy(conn) {
			return nil
		}
		_, err = g.reconnect(ctx, addr, conn)
		return err
	}
	conn, err = grpc.DialContext(ctx, addr, append(g.dopts, dopts...)...)
	if err != nil {
		return err
	}
	atomic.AddUint64(&g.clientCount, 1)
	g.conns.Store(addr, conn)
	return nil
}

func (g *gRPCClient) Disconnect(addr string) error {
	conn, ok := g.conns.Load(addr)
	if !ok {
		return errors.ErrGRPCClientConnNotFound(addr)
	}
	g.conns.Delete(addr)
	atomic.AddUint64(&g.clientCount, ^uint64(0))
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (g *gRPCClient) Close() error {
	closeList := make([]string, 0, int(atomic.LoadUint64(&g.clientCount)))
	g.conns.Range(func(addr string, conn *ClientConn) bool {
		if conn != nil {
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
	g.conns.Range(func(addr string, conn *ClientConn) bool {
		if g.isHealthy(conn) {
			connected = append(connected, addr)
		} else {
			disconnected = append(disconnected, addr)
		}
		return true
	})
	return
}

func (g *gRPCClient) isHealthy(conn *ClientConn) bool {
	return conn != nil &&
		conn.GetState() != connectivity.Shutdown &&
		conn.GetState() != connectivity.TransientFailure
}

func (g *gRPCClient) reconnect(ctx context.Context, addr string, conn *ClientConn) (rconn *ClientConn, err error) {
	defer func() {
		if err != nil {
			g.conns.Delete(addr)
		}
	}()
	if conn == nil {
		var ok bool
		conn, ok = g.conns.Load(addr)
		if !ok {
			return nil, errors.ErrGRPCClientConnNotFound(addr)
		}
	}
	if g.isHealthy(conn) {
		return conn, nil
	}
	if len(addr) != 0 {
		g.conns.Delete(addr)
	}
	if conn != nil {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
		conn = nil
	}
	conn, err = grpc.DialContext(ctx, addr, g.dopts...)
	if err != nil {
		return nil, err
	}
	g.conns.Store(addr, conn)
	return conn, nil
}
