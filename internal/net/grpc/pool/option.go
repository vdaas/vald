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

// Package pool provides grpc connection pool client
package pool

import "github.com/vdaas/vald/internal/backoff"

type Option func(*pool)

var (
	defaultOpts = []Option{
		WithSize(3),
		WithStartPort(80),
		WithEndPort(65535),
	}
)

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
		p.size = size
	}
}

func WithDialOptions(opts ...DialOption) Option {
	return func(p *pool) {
		if opts != nil && len(opts) > 0 {
			if len(p.dopts) > 0 {
				p.dopts = append(p.dopts, opts...)
			} else {
				p.dopts = opts
			}
		}
	}
}
