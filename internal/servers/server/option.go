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

// Package servers provides implementation of Go API for managing server flow
package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*server)

var (
	defaultOptions = []Option{
		WithDisableRestart(),
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
							log.Error(err, info.Get())
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
		if host != "" {
			s.host = host
		}
	}
}

func WithPort(port uint) Option {
	return func(s *server) {
		if port != 0 {
			s.port = port
		}
	}
}

func WithName(name string) Option {
	return func(s *server) {
		if name != "" {
			s.name = name
		}
	}
}

func WithErrorGroup(eg errgroup.Group) Option {
	return func(s *server) {
		if eg != nil {
			s.eg = eg
		}
	}
}

func WithPreStopFunction(f func() error) Option {
	return func(s *server) {
		if f != nil {
			s.preStopFunc = f
		}
	}
}

func WithPreStartFunc(f func() error) Option {
	return func(s *server) {
		if f != nil {
			s.preStartFunc = f
		}
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

func WithListenConfig(lc *net.ListenConfig) Option {
	return func(s *server) {
		if lc != nil {
			s.lc = lc
		}
	}
}

func WithServerMode(m ServerMode) Option {
	return func(s *server) {
		switch m {
		case GRPC, REST, GQL:
			s.mode = m
		}
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(s *server) {
		if cfg != nil {
			s.tcfg = cfg
		}
	}
}

func WithHTTPHandler(h http.Handler) Option {
	return func(s *server) {
		if h != nil {
			s.http.h = h
		}
	}
}

func WithHTTPServer(srv *http.Server) Option {
	return func(s *server) {
		if srv != nil {
			s.http.srv = srv
		}
	}
}

func WithGRPCServer(srv *grpc.Server) Option {
	return func(s *server) {
		if srv != nil {
			s.grpc.srv = srv
		}
	}
}

func WithGRPCOption(opts ...grpc.ServerOption) Option {
	return func(s *server) {
		if opts == nil {
			return
		}

		if s.grpc.opts == nil {
			s.grpc.opts = opts
			return
		}

		s.grpc.opts = append(s.grpc.opts, opts...)
	}
}

func WithGRPCRegistFunc(f func(*grpc.Server)) Option {
	return func(s *server) {
		if f != nil {
			s.grpc.reg = f
		}
	}
}

func WithEnableRestart() Option {
	return func(s *server) {
		s.enableRestart = true
	}
}

func WithDisableRestart() Option {
	return func(s *server) {
		s.enableRestart = false
	}
}

func WithGRPCMaxReceiveMessageSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxRecvMsgSize(size))
		}
	}
}

func WithGRPCMaxSendMessageSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxSendMsgSize(size))
		}
	}
}

func WithGRPCInitialWindowSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.InitialWindowSize(int32(size)))
		}
	}
}

func WithGRPCInitialConnWindowSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.InitialConnWindowSize(int32(size)))
		}
	}
}

func WithGRPCKeepaliveMaxConnIdle(max string) Option {
	return func(s *server) {
		if len(max) == 0 {
			return
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepAlive)
		}
		s.grpc.keepAlive.maxConnIdle = d
	}
}

func WithGRPCKeepaliveMaxConnAge(max string) Option {
	return func(s *server) {
		if len(max) == 0 {
			return
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepAlive)
		}
		s.grpc.keepAlive.maxConnAge = d
	}
}

func WithGRPCKeepaliveMaxConnAgeGrace(max string) Option {
	return func(s *server) {
		if len(max) == 0 {
			return
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepAlive)
		}
		s.grpc.keepAlive.maxConnAgeGrace = d
	}
}

func WithGRPCKeepaliveTime(dur string) Option {
	return func(s *server) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepAlive)
		}
		s.grpc.keepAlive.t = d
	}
}

func WithGRPCKeepaliveTimeout(dur string) Option {
	return func(s *server) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepAlive)
		}
		s.grpc.keepAlive.timeout = d
	}
}

func WithGRPCWriteBufferSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.WriteBufferSize(size))
		}
	}
}

func WithGRPCReadBufferSize(size int) Option {
	return func(s *server) {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.ReadBufferSize(size))
		}
	}
}

func WithGRPCConnectionTimeout(to string) Option {
	return func(s *server) {
		if len(to) == 0 {
			return
		}
		d, err := timeutil.Parse(to)
		if err != nil {
			return
		}
		s.grpc.opts = append(s.grpc.opts, grpc.ConnectionTimeout(d))
	}
}

func WithGRPCMaxHeaderListSize(size int) Option {
	return func(s *server) {
		if size > 0 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxHeaderListSize(uint32(size)))
		}
	}
}

func WithGRPCHeaderTableSize(size int) Option {
	return func(s *server) {
		if size > 0 {
			s.grpc.opts = append(s.grpc.opts, grpc.HeaderTableSize(uint32(size)))
		}
	}
}

func WithGRPCInterceptors(name ...string) Option {
	return func(s *server) {
		// s.grpc.opts = append(s.grpc.opts, grpc.UnaryInterceptor(uint32(size)))
	}
}
