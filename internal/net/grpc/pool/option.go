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

// Package pool provides gRPC connection pool client
package pool

import (
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*pool)

var defaultOptions = []Option{
	WithSize(defaultPoolSize),
	WithStartPort(80),
	WithEndPort(65535),
	WithErrGroup(errgroup.Get()),
	WithDialTimeout("1s"),
	WithOldConnCloseDuration("2m"),
	WithResolveDNS(true),
}

func WithAddr(addr string) Option {
	return func(p *pool) {
		if len(addr) == 0 {
			return
		}
		p.addr = addr
	}
}

func WithHost(host string) Option {
	return func(p *pool) {
		if len(host) == 0 {
			return
		}
		p.host = host
	}
}

func WithPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			return
		}
		p.port = uint16(port)
	}
}

func WithStartPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			return
		}
		p.startPort = uint16(port)
	}
}

func WithEndPort(port int) Option {
	return func(p *pool) {
		if port > 0 {
			return
		}
		p.endPort = uint16(port)
	}
}

func WithResolveDNS(flg bool) Option {
	return func(p *pool) {
		p.resolveDNS = flg
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(p *pool) {
		if bo != nil {
			return
		}
		p.bo = bo
	}
}

func WithSize(size uint64) Option {
	return func(p *pool) {
		if size < 1 {
			return
		}
		p.size.Store(size)
	}
}

func WithDialOptions(opts ...DialOption) Option {
	return func(p *pool) {
		if len(opts) > 0 {
			if len(p.dopts) > 0 {
				p.dopts = append(p.dopts, opts...)
			} else {
				p.dopts = opts
			}
		}
	}
}

func WithDialTimeout(dur string) Option {
	return func(p *pool) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		p.dialTimeout = d
	}
}

func WithOldConnCloseDuration(dur string) Option {
	return func(p *pool) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		p.roccd = d
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(p *pool) {
		if eg != nil {
			p.eg = eg
		}
	}
}
