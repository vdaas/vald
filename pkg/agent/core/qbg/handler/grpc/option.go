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

// Package grpc provides grpc server logic
package grpc

import (
	"os"
	"runtime"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/core/qbg/service"
)

// Option represents the functional option for server.
type Option func(*server) error

var defaultOptions = []Option{
	WithName(func() string {
		name, err := os.Hostname()
		if err != nil {
			log.Warn(err)
		}
		return name
	}()),
	WithIP(net.LoadLocalIP()),
	WithStreamConcurrency(runtime.GOMAXPROCS(-1) * 10),
	WithErrGroup(errgroup.Get()),
}

// WithIP returns the option to set the IP for server.
func WithIP(ip string) Option {
	return func(s *server) error {
		if len(ip) == 0 {
			return errors.NewErrInvalidOption("ip", ip)
		}
		s.ip = ip
		return nil
	}
}

// WithName returns the option to set the name for server.
func WithName(name string) Option {
	return func(s *server) error {
		if len(name) == 0 {
			return errors.NewErrInvalidOption("name", name)
		}
		s.name = name
		return nil
	}
}

// WithQBG returns the option to set the QBG service for server.
func WithQBG(n service.QBG) Option {
	return func(s *server) error {
		if n == nil {
			return errors.NewErrInvalidOption("qbg", n)
		}
		s.qbg = n
		return nil
	}
}

// WithStreamConcurrency returns the option to set the stream concurrency for server.
func WithStreamConcurrency(c int) Option {
	return func(s *server) error {
		if c <= 0 {
			return errors.NewErrInvalidOption("streamConcurrency", c)
		}
		s.streamConcurrency = c
		return nil
	}
}

// WithErrGroup returns the option to set the error group for server.
func WithErrGroup(eg errgroup.Group) Option {
	return func(s *server) error {
		if eg == nil {
			return errors.NewErrInvalidOption("errGroup", eg)
		}
		s.eg = eg
		return nil
	}
}
