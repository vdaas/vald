//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package discoverer
package discoverer

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(c *client) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithAutoConnect(true),
	WithNamespace("vald"),
}

func WithOnDiscoverFunc(f func(ctx context.Context, c Client, addrs []string) error) Option {
	return func(c *client) error {
		if f != nil {
			c.onDiscover = f
		}
		return nil
	}
}

func WithOnConnectFunc(f func(ctx context.Context, c Client, addr string) error) Option {
	return func(c *client) error {
		if f != nil {
			c.onConnect = f
		}
		return nil
	}
}

func WithOnDisconnectFunc(f func(ctx context.Context, c Client, addr string) error) Option {
	return func(c *client) error {
		if f != nil {
			c.onDisconnect = f
		}
		return nil
	}
}

func WithDiscovererClient(gc grpc.Client) Option {
	return func(c *client) error {
		c.dscClient = gc
		return nil
	}
}

func WithDiscoverDuration(dur string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		c.dscDur = d
		return nil
	}
}

func WithOptions(opts ...grpc.Option) Option {
	return func(c *client) error {
		c.opts = append(c.opts, opts...)
		return nil
	}
}

func WithAutoConnect(flg bool) Option {
	return func(c *client) error {
		c.autoconn = flg
		return nil
	}
}

func WithName(name string) Option {
	return func(c *client) error {
		if name != "" {
			c.name = name
		}
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(c *client) error {
		if ns != "" {
			c.namespace = ns
		}
		return nil
	}
}

func WithPort(port int) Option {
	return func(c *client) error {
		c.port = port
		return nil
	}
}

func WithServiceDNSARecord(a string) Option {
	return func(c *client) error {
		c.dns = a
		return nil
	}
}

func WithNodeName(nn string) Option {
	return func(c *client) error {
		if nn != "" {
			c.nodeName = nn
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}
