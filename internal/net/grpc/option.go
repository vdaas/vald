//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/client/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Option func(*gRPCClient)

var defaultOptions = []Option{
	WithConnectionPoolSize(3),
	WithEnableConnectionPoolRebalance(false),
	WithConnectionPoolRebalanceDuration("1h"),
	WithErrGroup(errgroup.Get()),
	WithHealthCheckDuration("10s"),
	WithResolveDNS(true),
	WithBackoffMaxDelay(gbackoff.DefaultConfig.MaxDelay.String()),
	WithBackoffBaseDelay(gbackoff.DefaultConfig.BaseDelay.String()),
	WithBackoffMultiplier(gbackoff.DefaultConfig.Multiplier),
	WithBackoffJitter(gbackoff.DefaultConfig.Jitter),
	WithMinConnectTimeout("20s"),
}

func WithAddrs(addrs ...string) Option {
	return func(g *gRPCClient) {
		if len(addrs) == 0 {
			return
		}
		if g.addrs == nil {
			g.addrs = make(map[string]struct{})
		}
		for _, addr := range addrs {
			g.addrs[addr] = struct{}{}
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

func WithConnectionPoolRebalanceDuration(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Hour
		}
		g.prDur = d
	}
}

func WithResolveDNS(flg bool) Option {
	return func(g *gRPCClient) {
		g.resolveDNS = flg
	}
}

func WithEnableConnectionPoolRebalance(flg bool) Option {
	return func(g *gRPCClient) {
		g.enablePoolRebalance = flg
	}
}

func WithConnectionPoolSize(size int) Option {
	return func(g *gRPCClient) {
		if size >= 1 {
			g.poolSize = uint64(size)
		}
	}
}

func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(g *gRPCClient) {
		if g.dopts != nil && len(g.dopts) > 0 {
			g.dopts = append(g.dopts, opts...)
		} else {
			g.dopts = opts
		}
	}
}

func WithBackoffMaxDelay(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.gbo.MaxDelay = d
	}
}

func WithBackoffBaseDelay(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.gbo.BaseDelay = d
	}
}

func WithBackoffMultiplier(m float64) Option {
	return func(g *gRPCClient) {
		g.gbo.Multiplier = m
	}
}

func WithBackoffJitter(j float64) Option {
	return func(g *gRPCClient) {
		g.gbo.Jitter = j
	}
}

func WithMinConnectTimeout(dur string) Option {
	return func(g *gRPCClient) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.mcd = d
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

func WithCircuitBreaker(cb circuitbreaker.CircuitBreaker) Option {
	return func(gr *gRPCClient) {
		if cb != nil {
			gr.cb = cb
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
			g.dopts = append(g.dopts,
				grpc.WithWriteBufferSize(size),
			)
		}
	}
}

func WithReadBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.dopts = append(g.dopts,
				grpc.WithReadBufferSize(size),
			)
		}
	}
}

func WithInitialWindowSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.dopts = append(g.dopts,
				grpc.WithInitialWindowSize(int32(size)),
			)
		}
	}
}

func WithInitialConnectionWindowSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.dopts = append(g.dopts,
				grpc.WithInitialConnWindowSize(int32(size)),
			)
		}
	}
}

func WithMaxMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.dopts = append(g.dopts,
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(size)),
			)
		}
	}
}

func WithInsecure(flg bool) Option {
	return func(g *gRPCClient) {
		if flg {
			g.dopts = append(g.dopts,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
		}
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
		g.dopts = append(g.dopts,
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

func WithDialer(der net.Dialer) Option {
	return func(g *gRPCClient) {
		if der != nil {
			g.dialer = der
			g.dopts = append(g.dopts,
				grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
					// TODO we need change network type dynamically
					log.Debugf("grpc Context Dialer addr is %s", addr)
					return der.GetDialer()(ctx, net.TCP.String(), addr)
				}),
			)
		}
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(g *gRPCClient) {
		if cfg != nil {
			g.dopts = append(g.dopts,
				grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
			)
		}
	}
}

func WithClientInterceptors(names ...string) Option {
	return func(g *gRPCClient) {
		for _, name := range names {
			switch strings.ToLower(name) {
			case "traceinterceptor", "trace":
				g.dopts = append(g.dopts,
					grpc.WithUnaryInterceptor(trace.UnaryClientInterceptor()),
					grpc.WithStreamInterceptor(trace.StreamClientInterceptor()),
				)
			default:
			}
		}
	}
}

func WithOldConnCloseDuration(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		g.roccd = dur
	}
}
