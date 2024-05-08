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

// Package starter provides server startup and shutdown flow control
package starter

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/servers/server"
)

type Option func(*srvs)

func WithConfig(cfg *config.Servers) Option {
	return func(s *srvs) {
		s.cfg = cfg
	}
}

func WithGRPC(opts func(cfg *config.Server) []server.Option) Option {
	return func(s *srvs) {
		s.grpc = opts
	}
}

func WithREST(opts func(cfg *config.Server) []server.Option) Option {
	return func(s *srvs) {
		s.rest = opts
	}
}

func WithGQL(opts func(cfg *config.Server) []server.Option) Option {
	return func(s *srvs) {
		s.gql = opts
	}
}

func WithPreStartFunc(name string, f func() error) Option {
	return func(s *srvs) {
		if f != nil && s.pstartf != nil {
			s.pstartf[name] = f
		}
	}
}

func WithPreStopFunc(name string, f func() error) Option {
	return func(s *srvs) {
		if f != nil && s.pstopf != nil {
			s.pstopf[name] = f
		}
	}
}
