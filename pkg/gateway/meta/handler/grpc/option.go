//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package grpc

import (
	"runtime"

	"github.com/vdaas/vald/internal/client/v1/client/meta"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Option func(*server)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithStreamConcurrency(runtime.GOMAXPROCS(-1) * 10),
	WithMultiConcurrency(runtime.GOMAXPROCS(-1) * 10),
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

func WithMetadataClient(mc meta.MetadataClient) Option {
	return func(s *server) {
		if mc != nil {
			s.metadataClient = mc
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

func WithMultiConcurrency(c int) Option {
	return func(s *server) {
		if c > 1 {
			s.multiConcurrency = c
		}
	}
}
