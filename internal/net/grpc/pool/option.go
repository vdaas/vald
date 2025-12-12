//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package pool

import (
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

// Option defines a functional option for configuring the pool.
type Option func(*pool)

// Default options.
var defaultOptions = []Option{
	WithSize(defaultPoolSize),
	WithStartPort(80),
	WithEndPort(65535),
	WithErrGroup(errgroup.Get()),
	WithDialTimeout("1s"),
	WithOldConnCloseDelay("20s"),
	WithResolveDNS(true),
}

// WithAddr sets the target address. It also extracts the host and port.
func WithAddr(addr string) Option {
	return func(p *pool) {
		if addr == "" {
			return
		}
		p.addr = addr
		var err error
		// Attempt to split host and port.
		if p.host, p.port, err = net.SplitHostPort(addr); err != nil {
			p.host = addr
		}
	}
}

// WithHost sets the target host.
func WithHost(host string) Option {
	return func(p *pool) {
		if host != "" {
			p.host = host
		}
	}
}

// WithPort sets the target port.
func WithPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			p.port = uint16(port)
		}
	}
}

// WithStartPort sets the starting port for scanning.
func WithStartPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			p.startPort = uint16(port)
		}
	}
}

// WithEndPort sets the ending port for scanning.
func WithEndPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			p.endPort = uint16(port)
		}
	}
}

// WithResolveDNS enables or disables DNS resolution.
func WithResolveDNS(enable bool) Option {
	return func(p *pool) {
		p.enableDNSLookup = enable
	}
}

// WithBackoff sets the backoff strategy.
func WithBackoff(bo backoff.Backoff) Option {
	return func(p *pool) {
		if bo != nil {
			p.bo = bo
		}
	}
}

// WithSize sets the pool size.
func WithSize(size uint64) Option {
	return func(p *pool) {
		if size < 1 {
			return
		}
		p.poolSize.Store(size)
	}
}

// WithDialOptions appends gRPC dial options.
func WithDialOptions(opts ...DialOption) Option {
	return func(p *pool) {
		if len(opts) > 0 {
			p.dialOpts = append(p.dialOpts, opts...)
		}
	}
}

// WithDialTimeout sets the dial timeout duration.
func WithDialTimeout(dur string) Option {
	return func(p *pool) {
		if dur == "" {
			return
		}
		if t, err := timeutil.Parse(dur); err == nil {
			p.dialTimeout = t
		}
	}
}

// WithOldConnCloseDelay sets the delay before closing old connections.
func WithOldConnCloseDelay(dur string) Option {
	return func(p *pool) {
		if dur == "" {
			return
		}
		if t, err := timeutil.Parse(dur); err == nil {
			p.oldConnCloseDelay = t
		}
	}
}

// WithErrGroup sets the errgroup for goroutine management.
func WithErrGroup(eg errgroup.Group) Option {
	return func(p *pool) {
		if eg != nil {
			p.errGroup = eg
		}
	}
}
