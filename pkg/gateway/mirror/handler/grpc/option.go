// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"os"
	"runtime"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type Option func(*server) error

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
	return func(s *server) error {
		if len(ip) != 0 {
			s.ip = ip
		}
		return nil
	}
}

// WithName returns the option to set the name for server.
func WithName(name string) Option {
	return func(s *server) error {
		if len(name) != 0 {
			s.name = name
		}
		return nil
	}
}

func WithGateway(g service.Gateway) Option {
	return func(s *server) error {
		if g != nil {
			s.gateway = g
		}
		return nil
	}
}

func WithMirror(m service.Mirror) Option {
	return func(s *server) error {
		if m != nil {
			s.mirror = m
		}
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(s *server) error {
		if eg != nil {
			s.eg = eg
		}
		return nil
	}
}

func WithStreamConcurrency(c int) Option {
	return func(s *server) error {
		if c > 0 {
			s.streamConcurrency = c
		}
		return nil
	}
}

func WithValdAddr(addr string) Option {
	return func(s *server) error {
		if len(addr) == 0 {
			return errors.NewErrCriticalOption("valdAddr", addr)
		}
		s.vAddr = addr
		return nil
	}
}
