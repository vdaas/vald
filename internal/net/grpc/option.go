//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
	"crypto/tls"
	"net"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Option func(*gRPCClient)

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithHealthCheckDuration("10s"),
	}
)

func WithAddrs(addrs ...string) Option {
	return func(g *gRPCClient) {
		if len(addrs) == 0 {
			return
		}
		if g.addrs == nil || len(g.addrs) == 0 {
			g.addrs = addrs
		} else {
			g.addrs = append(g.addrs, addrs...)
		}
	}
}

func WithHealthCheckDuration(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.hcDur = d
	}
}

func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(g *gRPCClient) {
		if g.gopts != nil && len(g.gopts) > 0 {
			g.gopts = append(g.gopts, opts...)
		} else {
			g.gopts = opts
		}
	}
}

func WithMaxBackoffDelay(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.gopts = append(g.gopts,
			grpc.WithBackoffMaxDelay(d),
		)
	}
}

func WithCallOptions(opts ...grpc.CallOption) Option {
	return func(g *gRPCClient) {
		if g.copts != nil && len(g.copts) > 0 {
			g.copts = append(g.copts, opts...)
		} else {
			g.copts = opts
		}
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gRPCClient) {
		if eg != nil {
			g.eg = eg
		}
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(g *gRPCClient) {
		if bo != nil {
			g.bo = bo
		}
	}
}

func WithWaitForReady(flg bool) Option {
	return func(g *gRPCClient) {
		g.copts = append(g.copts,
			grpc.WaitForReady(flg),
		)
	}
}
func WithMaxRetryRPCBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.copts = append(g.copts,
				grpc.MaxRetryRPCBufferSize(size),
			)
		}
	}
}
func WithMaxRecvMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.copts = append(g.copts,
				grpc.MaxCallRecvMsgSize(size),
			)
		}
	}
}
func WithMaxSendMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.copts = append(g.copts,
				grpc.MaxCallSendMsgSize(size),
			)
		}
	}
}
func WithWriteBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.gopts = append(g.gopts,
				grpc.WithWriteBufferSize(size),
			)
		}
	}
}
func WithReadBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.gopts = append(g.gopts,
				grpc.WithReadBufferSize(size),
			)
		}
	}
}
func WithInitialWindowSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.gopts = append(g.gopts,
				grpc.WithInitialWindowSize(int32(size)),
			)
		}
	}
}
func WithInitialConnectionWindowSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.gopts = append(g.gopts,
				grpc.WithInitialConnWindowSize(int32(size)),
			)
		}
	}
}
func WithMaxMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.gopts = append(g.gopts,
				grpc.WithMaxMsgSize(size),
			)
		}
	}
}

func WithInsecure(flg bool) Option {
	return func(g *gRPCClient) {
		if flg {
			g.gopts = append(g.gopts,
				grpc.WithInsecure(),
			)
		}
	}
}

func WithDialTimeout(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		g.gopts = append(g.gopts,
			grpc.WithTimeout(d),
		)
	}
}

func WithKeepaliveParams(t, to string, permitWithoutStream bool) Option {
	return func(g *gRPCClient) {
		if len(t) == 0 || len(to) == 0 {
			return
		}
		td, err := timeutil.Parse(t)
		if err != nil {
			return
		}
		tod, err := timeutil.Parse(t)
		if err != nil {
			return
		}
		g.gopts = append(g.gopts,
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                td,
					Timeout:             tod,
					PermitWithoutStream: permitWithoutStream,
				},
			),
		)
	}
}

func WithDialer(der tcp.Dialer) Option {
	return func(g *gRPCClient) {
		if der != nil {
			g.gopts = append(g.gopts,
				grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
					return der.GetDialer()(ctx, "tcp", addr)
				}),
			)
		}
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(g *gRPCClient) {
		if cfg != nil {
			g.gopts = append(g.gopts,
				grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
			)
		}
	}
}
