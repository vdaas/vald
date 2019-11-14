//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package grpc provides generic functionallity for grpc
package grpc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientConn = grpc.ClientConn
type CallOption = grpc.CallOption
type DialOption = grpc.DialOption

type Client interface {
	StartConnectionMonitor(ctx context.Context) <-chan error
	Connect(ctx context.Context, addr string) error
	Disconnect(addr string) error
	Range(ctx context.Context,
		f func(addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption) error) error
	RangeConcurrent(ctx context.Context,
		concurrency int,
		f func(addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption) error) error
	Do(ctx context.Context,
		addr string, f func(conn *grpc.ClientConn,
			copts ...grpc.CallOption) (interface{}, error)) (interface{}, error)
	GetDialOption() []grpc.DialOption
	GetCallOption() []grpc.CallOption
	Close() error
}

type gRPCClient struct {
	fch   chan string
	addrs []string
	conns gRPCConns
	hcDur time.Duration
	gopts []grpc.DialOption
	copts []grpc.CallOption
	eg    errgroup.Group
	bo    backoff.Backoff
}

func New(opts ...Option) (c Client) {
	g := new(gRPCClient)
	for _, opt := range append(defaultOpts, opts...) {
		opt(g)
	}

	g.fch = make(chan string, len(g.addrs))

	return g
}

func (g *gRPCClient) StartConnectionMonitor(ctx context.Context) <-chan error {
	ech := make(chan error, len(g.addrs))
	if g.addrs == nil || len(g.addrs) == 0 {
		defer close(ech)
		return ech
	}

	for _, addr := range g.addrs {
		conn, err := grpc.DialContext(ctx, addr,
			append(g.gopts, grpc.WithBlock())...)
		if err != nil {
			ech <- err
		} else {
			g.conns.Store(addr, conn)
		}
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
			case addr := <-g.fch:
				err = g.Connect(ctx, addr)
				if err != nil {
					ech <- err
				}
			case <-tick.C:
				g.conns.Range(func(addr string, conn *grpc.ClientConn) bool {
					if conn == nil ||
						conn.GetState() == connectivity.Shutdown ||
						conn.GetState() == connectivity.TransientFailure {
						err = g.Disconnect(addr)
						if err != nil {
							ech <- err
						}
						err = g.Connect(ctx, addr)
						if err != nil {
							ech <- err
						}
					}
					return true
				})
			}
		}
	}))
	return ech
}

func (g *gRPCClient) Range(ctx context.Context,
	f func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error) (rerr error) {
	g.conns.Range(func(addr string, conn *grpc.ClientConn) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			var err error
			if g.bo != nil {
				_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
					err = f(addr, conn, g.copts...)
					if err != nil {
						err = g.Connect(ctx, addr)
					}
					return
				})
			} else {
				err = f(addr, conn, g.copts...)
			}

			if err != nil {
				rerr = errors.Wrap(rerr, fmt.Sprintf("addr: %s\terror: %s", addr, err.Error()))
			}
		}
		return true
	})
	return rerr
}

func (g *gRPCClient) RangeConcurrent(ctx context.Context,
	concurrency int,
	f func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error) (rerr error) {
	eg, ctx := errgroup.New(ctx)
	eg.Limitation(concurrency)
	g.conns.Range(func(addr string, conn *grpc.ClientConn) bool {
		eg.Go(func() (err error) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				var err error
				if g.bo != nil {
					_, err = g.bo.Do(ctx, func() (r interface{}, err error) {
						err = f(addr, conn, g.copts...)
						if err != nil {
							err = g.Connect(ctx, addr)
						}
						return
					})
				} else {
					err = f(addr, conn, g.copts...)
				}
				if err != nil {
					return errors.Wrap(rerr, fmt.Sprintf("addr: %s\terror: %s", addr, err.Error()))
				}
			}
			return nil
		})
		return true
	})
	return eg.Wait()
}

func (g *gRPCClient) Do(ctx context.Context, addr string,
	f func(conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error)) (data interface{}, err error) {
	conn, ok := g.conns.Load(addr)
	if !ok {
		return nil, errors.ErrGRPCClientConnNotFound(addr)
	}
	if g.bo != nil {
		data, err = g.bo.Do(ctx, func() (r interface{}, err error) {
			r, err = f(conn, g.copts...)
			if err != nil {
				err = g.Connect(ctx, addr)
				return nil, err
			}
			return r, nil
		})
	} else {
		data, err = f(conn, g.copts...)
	}
	if err != nil {
		return nil, errors.Errorf("addr: %s\terror: %s", addr, err.Error())
	}
	return
}

func (g *gRPCClient) GetDialOption() []grpc.DialOption {
	return g.gopts
}

func (g *gRPCClient) GetCallOption() []grpc.CallOption {
	return g.copts
}

func (g *gRPCClient) Connect(ctx context.Context, addr string) error {
	conn, ok := g.conns.Load(addr)
	if ok {
		if conn == nil ||
			conn.GetState() == connectivity.Shutdown ||
			conn.GetState() == connectivity.TransientFailure {
			g.Disconnect(addr)
		} else {
			return nil
		}
	}
	conn, err := grpc.DialContext(ctx, addr, g.gopts...)
	if err != nil {
		go func() {
			time.Sleep(time.Second)
			g.fch <- addr
		}()
		runtime.Gosched()
		return err
	}
	g.conns.Store(addr, conn)
	return nil
}

func (g *gRPCClient) Disconnect(addr string) error {
	conn, ok := g.conns.Load(addr)
	if !ok {
		return errors.ErrGRPCClientConnNotFound(addr)
	}
	g.conns.Delete(addr)
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (g *gRPCClient) Close() error {
	g.conns.Range(func(addr string, conn *grpc.ClientConn) bool {
		if conn != nil {
			g.Disconnect(addr)
		}
		return true
	})
	close(g.fch)
	return nil
}
