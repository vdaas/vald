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

// Package service
package service

import (
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
)

type GWOption func(g *gateway) error

var (
	defaultGWOpts = []GWOption{}
)

func WithDiscoverHost(host string) GWOption {
	return func(g *gateway) error {
		g.dscHost = host
		return nil
	}
}

func WithDiscoverPort(port int) GWOption {
	return func(g *gateway) error {
		g.dscPort = port
		return nil
	}
}

func WithDiscoverDuration(dur string) GWOption {
	return func(g *gateway) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		g.dscDur = d
		return nil
	}
}

func WithPort(port int) GWOption {
	return func(g *gateway) error {
		g.port = port
		return nil
	}
}

func WithGRPCDialOption(opt grpc.DialOption) GWOption {
	return func(g *gateway) error {
		g.gopts = append(g.gopts, opt)
		return nil
	}
}

func WithGRPCDialOptions(opts []grpc.DialOption) GWOption {
	return func(g *gateway) error {
		if g.gopts != nil && len(g.gopts) > 0 {
			g.gopts = append(g.gopts, opts...)
		} else {
			g.gopts = opts
		}
		return nil
	}
}

func withBackoff(bo backoff.Backoff) GWOption {
	return func(g *gateway) error {
		g.bo = bo
		return nil
	}
}

func withErrGroup(eg errgroup.Group) GWOption {
	return func(g *gateway) error {
		g.eg = eg
		return nil
	}
}
