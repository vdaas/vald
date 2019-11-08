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
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
)

type Option func(*gRPCClient) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithHealthCheckDuration("10s"),
	}
)

func WithAddrs(addrs ...string) Option {
	return func(g *gRPCClient) error {
		if g.addrs == nil || len(g.addrs) == 0 {
			g.addrs = addrs
		} else {
			g.addrs = append(g.addrs, addrs...)
		}
		return nil
	}
}

func WithHealthCheckDuration(dur string) Option {
	return func(g *gRPCClient) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.hcDur = d
		return nil
	}
}

func WithGRPCDialOptions(opts ...grpc.DialOption) Option {
	return func(g *gRPCClient) error {
		if g.gopts != nil && len(g.gopts) > 0 {
			g.gopts = append(g.gopts, opts...)
		} else {
			g.gopts = opts
		}
		return nil
	}
}

func WithGRPCCallOptions(opts ...grpc.CallOption) Option {
	return func(g *gRPCClient) error {
		if g.copts != nil && len(g.copts) > 0 {
			g.copts = append(g.copts, opts...)
		} else {
			g.copts = opts
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gRPCClient) error {
		if eg != nil {
			g.eg = eg
		}
		return nil
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(g *gRPCClient) error {
		if bo != nil {
			g.bo = bo
		}
		return nil
	}
}
