//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc server logic
package grpc

import (
	"os"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/pkg/gateway/meta/service"
)

type Option func(*server)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithStreamConcurrency(20),
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

func WithValdClient(g vald.Client) Option {
	return func(s *server) {
		if g != nil {
			s.gateway = g
		}
	}
}

func WithMeta(m service.Meta) Option {
	return func(s *server) {
		if m != nil {
			s.metadata = m
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
		if c != 0 {
			s.streamConcurrency = c
		}
	}
}
