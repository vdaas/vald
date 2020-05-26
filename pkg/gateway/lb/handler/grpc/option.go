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

// Package grpc provides grpc server logic
package grpc

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Option func(*server)

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithReplicationCount(3),
		WithStreamConcurrency(20),
		WithTimeout("5s"),
	}
)

func WithGateway(g service.Gateway) Option {
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

func WithBackup(b service.Backup) Option {
	return func(s *server) {
		if b != nil {
			s.backup = b
		}
	}
}

func WithFilters(filter service.Filter) Option {
	return func(s *server) {
		if filter != nil {
			s.filter = filter
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

func WithTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 10
		}
		s.timeout = d
	}
}

func WithReplicationCount(rep int) Option {
	return func(s *server) {
		if rep > 1 {
			s.replica = rep
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
