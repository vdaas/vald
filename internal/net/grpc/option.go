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
	"net"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
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
		if g.addrs == nil || len(g.addrs) == 0 {
			g.addrs = addrs
		} else {
			g.addrs = append(g.addrs, addrs...)
		}
		return
	}
}

func WithHealthCheckDuration(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.hcDur = d
		return
	}
}

func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(g *gRPCClient) {
		if g.gopts != nil && len(g.gopts) > 0 {
			g.gopts = append(g.gopts, opts...)
		} else {
			g.gopts = opts
		}
		return
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
		return
	}
}

func WithCallOptions(opts ...grpc.CallOption) Option {
	return func(g *gRPCClient) {
		if g.copts != nil && len(g.copts) > 0 {
			g.copts = append(g.copts, opts...)
		} else {
			g.copts = opts
		}
		return
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gRPCClient) {
		if eg != nil {
			g.eg = eg
		}
		return
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(g *gRPCClient) {
		if bo != nil {
			g.bo = bo
		}
		return
	}
}

func WithWaitForReady(flg bool) Option {
	return func(g *gRPCClient) {
		g.copts = append(g.copts,
			grpc.WaitForReady(flg),
		)
		return
	}
}
func WithMaxRetryRPCBufferSize(size int) Option {
	return func(g *gRPCClient) {
		g.copts = append(g.copts,
			grpc.MaxRetryRPCBufferSize(size),
		)
		return
	}
}
func WithMaxRecvMsgSize(size int) Option {
	return func(g *gRPCClient) {
		g.copts = append(g.copts,
			grpc.MaxCallRecvMsgSize(size),
		)
		return
	}
}
func WithMaxSendMsgSize(size int) Option {
	return func(g *gRPCClient) {
		g.copts = append(g.copts,
			grpc.MaxCallSendMsgSize(size),
		)
		return
	}
}
func WithWriteBufferSize(size int) Option {
	return func(g *gRPCClient) {
		g.gopts = append(g.gopts,
			grpc.WithWriteBufferSize(size),
		)
		return
	}
}
func WithReadBufferSize(size int) Option {
	return func(g *gRPCClient) {
		g.gopts = append(g.gopts,
			grpc.WithReadBufferSize(size),
		)
		return
	}
}
func WithInitialWindowSize(size int) Option {
	return func(g *gRPCClient) {
		g.gopts = append(g.gopts,
			grpc.WithInitialWindowSize(int32(size)),
		)
		return
	}
}
func WithInitialConnectionWindowSize(size int) Option {
	return func(g *gRPCClient) {
		g.gopts = append(g.gopts,
			grpc.WithInitialConnWindowSize(int32(size)),
		)
		return
	}
}
func WithMaxMsgSize(size int) Option {
	return func(g *gRPCClient) {
		g.gopts = append(g.gopts,
			grpc.WithMaxMsgSize(size),
		)
		return
	}
}

func WithInsecure(flg bool) Option {
	return func(g *gRPCClient) {
		if flg {
			g.gopts = append(g.gopts,
				grpc.WithInsecure(),
			)
		}
		return
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
		return
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
		return
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
		return
	}
}
