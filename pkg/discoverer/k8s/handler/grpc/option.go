//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

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
}

func WithDiscoverer(dsc service.Discoverer) Option {
	return func(s *server) error {
		if dsc != nil {
			s.dsc = dsc
		}
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
