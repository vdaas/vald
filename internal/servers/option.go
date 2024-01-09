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

// Package servers provides implementation of Go API for managing server flow
package servers

import (
	"time"

	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*listener)

var defaultOptions = []Option{
	WithErrorGroup(errgroup.Get()),
}

func WithServer(srv server.Server) Option {
	return func(l *listener) {
		if srv == nil {
			return
		}
		if l.servers == nil {
			l.servers = make(map[string]server.Server)
		}
		l.servers[srv.Name()] = srv
	}
}

func WithErrorGroup(eg errgroup.Group) Option {
	return func(l *listener) {
		l.eg = eg
	}
}

func WithShutdownDuration(dur string) Option {
	return func(l *listener) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 20
		}
		l.sddur = d
	}
}

func WithStartUpStrategy(strg []string) Option {
	return func(l *listener) {
		l.sus = strg
	}
}

func WithShutdownStrategy(strg []string) Option {
	return func(l *listener) {
		l.sds = strg
	}
}
