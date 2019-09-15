//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package servers provides implementation of Go API for managing server flow
package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
)

type Option func(*server)

var (
	defaultOpts = []Option{
		WithServerMode(REST),
		WithErrorGroup(errgroup.Get()),
	}
	HealthServerOpts = func(name, host, path string, port uint) []Option {
		return []Option{
			WithName(name),
			WithErrorGroup(errgroup.Get()),
			WithHTTPHandler(func() http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet {
						w.Header().Set(rest.ContentType, rest.TextPlain+";"+rest.CharsetUTF8)
						w.WriteHeader(http.StatusOK)
						_, err := fmt.Fprint(w, http.StatusText(http.StatusOK))
						if err != nil {
							log.Fatal(err)
						}
					}
				})
				return mux
			}()),
			WithHost(host),
			WithIdleTimeout("3s"),
			WithPort(port),
			WithProbeWaitTime("2s"),
			WithReadHeaderTimeout("3s"),
			WithReadTimeout("2s"),
			WithServerMode(REST),
			WithShutdownDuration("4s"),
			WithWriteTimeout("3s"),
		}
	}
)

func WithHost(host string) Option {
	return func(s *server) {
		if host == "" {
			return
		}
		s.host = host
	}
}

func WithPort(port uint) Option {
	return func(s *server) {
		if port == 0 {
			return
		}
		s.port = port
	}
}

func WithName(name string) Option {
	return func(s *server) {
		s.name = name
	}
}

func WithErrorGroup(eg errgroup.Group) Option {
	return func(s *server) {
		s.eg = eg
	}
}

func WithPreStopFunction(f func() error) Option {
	return func(s *server) {
		s.preStopFunc = f
	}
}

func WithPreStartFunc(f func() error) Option {
	return func(s *server) {
		s.preStartFunc = f
	}
}

func WithProbeWaitTime(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.pwt = d
	}
}

func WithShutdownDuration(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.sddur = d
	}
}

func WithReadHeaderTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.rht = d
	}
}

func WithReadTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.rt = d
	}
}

func WithWriteTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.wt = d
	}
}

func WithIdleTimeout(dur string) Option {
	return func(s *server) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.it = d
	}
}

func WithListener(l net.Listener) Option {
	return func(s *server) {
		if l != nil {
			s.l = l
		}
	}
}

func WithServerMode(m mode) Option {
	return func(s *server) {
		s.mode = m
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(s *server) {
		if cfg != nil {
			return
		}
		s.tcfg = cfg
	}
}

func WithHTTPHandler(h http.Handler) Option {
	return func(s *server) {
		s.http.h = h
	}
}

func WithHTTPServer(srv *http.Server) Option {
	return func(s *server) {
		s.http.srv = srv
	}
}

func WithGRPCServer(srv *grpc.Server) Option {
	return func(s *server) {
		s.grpc.srv = srv
	}
}

func WithGRPCOption(opts ...grpc.ServerOption) Option {
	return func(s *server) {
		s.grpc.opts = opts
	}
}

func WithGRPCRegistFunc(f func(*grpc.Server)) Option {
	return func(s *server) {
		s.grpc.reg = f
	}
}
