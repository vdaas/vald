//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc server logic
package grpc

import (
	"os"
	"runtime"

	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Option func(*server)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithStreamConcurrency(runtime.GOMAXPROCS(-1) * 10),
	WithName(func() string {
		name, err := os.Hostname()
		if err != nil {
			log.Warn(err)
		}
		return name
	}()),
	WithIP(net.LoadLocalIP()),
}

// WithIP returns the option to set the IP for server.
func WithIP(ip string) Option {
	return func(s *server) {
		if len(ip) != 0 {
			s.ip = ip
		}
	}
}

// WithName returns the option to set the name for server.
func WithName(name string) Option {
	return func(s *server) {
		if len(name) != 0 {
			s.name = name
		}
	}
}

func WithIngressFilterClient(c ingress.Client) Option {
	return func(s *server) {
		if c != nil {
			s.ingress = c
		}
	}
}

func WithEgressFilterClient(c egress.Client) Option {
	return func(s *server) {
		if c != nil {
			s.egress = c
		}
	}
}

func WithValdClient(g vald.Client) Option {
	return func(s *server) {
		if g != nil {
			s.gateway = g
		}
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(s *server) {
		if eg != nil {
			s.eg = eg
		}
	}
}

func WithStreamConcurrency(c int) Option {
	return func(s *server) {
		if c > 1 {
			s.streamConcurrency = c
		}
	}
}

func WithVectorizerTargets(addr string) Option {
	return func(s *server) {
		if len(addr) == 0 {
			return
		}
		s.Vectorizer = addr
	}
}

func WithDistanceFilterTargets(cs ...*config.DistanceFilterConfig) Option {
	return func(s *server) {
		if len(cs) == 0 {
			return
		}
		if len(s.DistanceFilters) == 0 {
			s.DistanceFilters = cs
		} else {
			s.DistanceFilters = append(s.DistanceFilters, cs...)
		}
	}
}

func WithObjectFilterTargets(cs ...*config.ObjectFilterConfig) Option {
	return func(s *server) {
		if len(cs) == 0 {
			return
		}
		if len(s.ObjectFilters) == 0 {
			s.ObjectFilters = cs
		} else {
			s.ObjectFilters = append(s.ObjectFilters, cs...)
		}
	}
}

func WithSearchFilterTargets(addrs ...string) Option {
	return func(s *server) {
		if len(addrs) == 0 {
			return
		}
		if len(s.SearchFilters) == 0 {
			s.SearchFilters = addrs
		} else {
			s.SearchFilters = append(s.SearchFilters, addrs...)
		}
	}
}

func WithInsertFilterTargets(addrs ...string) Option {
	return func(s *server) {
		if len(addrs) == 0 {
			return
		}
		if len(s.InsertFilters) == 0 {
			s.InsertFilters = addrs
		} else {
			s.InsertFilters = append(s.InsertFilters, addrs...)
		}
	}
}

func WithUpdateFilterTargets(addrs ...string) Option {
	return func(s *server) {
		if len(addrs) == 0 {
			return
		}
		if len(s.UpdateFilters) == 0 {
			s.UpdateFilters = addrs
		} else {
			s.UpdateFilters = append(s.UpdateFilters, addrs...)
		}
	}
}

func WithUpsertFilterTargets(addrs ...string) Option {
	return func(s *server) {
		if len(addrs) == 0 {
			return
		}
		if len(s.UpsertFilters) == 0 {
			s.UpsertFilters = addrs
		} else {
			s.UpsertFilters = append(s.UpsertFilters, addrs...)
		}
	}
}
