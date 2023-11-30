//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/logging"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/metric"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/recover"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/server/trace"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*server) error

var (
	defaultOptions = []Option{
		WithDisableRestart(),
		WithNetwork(net.TCP.String()),
		WithServerMode(REST),
		WithErrorGroup(errgroup.Get()),
	}
	HealthServerOpts = func(name, host, path string, port uint16) []Option {
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
			WithNetwork(net.TCP.String()),
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

func WithNetwork(network string) Option {
	return func(s *server) error {
		if network != "" {
			nt := net.NetworkTypeFromString(network)
			if nt == 0 || nt == net.Unknown || strings.EqualFold(nt.String(), net.Unknown.String()) {
				nt = net.TCP
			}
			s.network = nt
		}
		return nil
	}
}

func WithSocketPath(path string) Option {
	return func(s *server) error {
		if path != "" {
			s.socketPath = path
		}
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *server) error {
		if host != "" {
			s.host = host
		}
		return nil
	}
}

func WithPort(port uint16) Option {
	return func(s *server) error {
		if port != 0 {
			s.port = port
		}
		return nil
	}
}

func WithSocketFlag(flg control.SocketFlag) Option {
	return func(s *server) error {
		s.sockFlg = flg
		return nil
	}
}

func WithName(name string) Option {
	return func(s *server) error {
		if name != "" {
			s.name = name
		}
		return nil
	}
}

func WithErrorGroup(eg errgroup.Group) Option {
	return func(s *server) error {
		if eg != nil {
			s.eg = eg
		}
		return nil
	}
}

func WithPreStopFunction(f func() error) Option {
	return func(s *server) error {
		if f != nil {
			s.preStopFunc = f
		}
		return nil
	}
}

func WithPreStartFunc(f func() error) Option {
	return func(s *server) error {
		if f != nil {
			s.preStartFunc = f
		}
		return nil
	}
}

func WithProbeWaitTime(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.pwt = d
		return nil
	}
}

func WithShutdownDuration(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.sddur = d
		return nil
	}
}

func WithReadHeaderTimeout(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.rht = d
		return nil
	}
}

func WithReadTimeout(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.rt = d
		return nil
	}
}

func WithWriteTimeout(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.wt = d
		return nil
	}
}

func WithIdleTimeout(dur string) Option {
	return func(s *server) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 5
		}
		s.it = d
		return nil
	}
}

func WithListenConfig(lc *net.ListenConfig) Option {
	return func(s *server) error {
		if lc != nil {
			s.lc = lc
		}
		return nil
	}
}

func WithServerMode(m ServerMode) Option {
	return func(s *server) error {
		switch m {
		case GRPC, REST, GQL:
			s.mode = m
		}
		return nil
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(s *server) error {
		if cfg != nil {
			s.tcfg = cfg
		}
		return nil
	}
}

func WithHTTPHandler(h http.Handler) Option {
	return func(s *server) error {
		if h != nil {
			s.http.h = h
		}
		return nil
	}
}

func WithHTTPServer(srv *http.Server) Option {
	return func(s *server) error {
		if srv != nil {
			s.http.srv = srv
		}
		return nil
	}
}

func WithGRPCServer(srv *grpc.Server) Option {
	return func(s *server) error {
		if srv != nil {
			s.grpc.srv = srv
		}
		return nil
	}
}

func WithGRPCOption(opts ...grpc.ServerOption) Option {
	return func(s *server) error {
		if opts == nil {
			return nil
		}

		if s.grpc.opts == nil {
			s.grpc.opts = opts
			return nil
		}

		s.grpc.opts = append(s.grpc.opts, opts...)
		return nil
	}
}

func WithGRPCRegistFunc(f func(*grpc.Server)) Option {
	return func(s *server) error {
		if f != nil {
			if s.grpc.regs == nil {
				s.grpc.regs = make([]func(*grpc.Server), 0, 2)
			}
			s.grpc.regs = append(s.grpc.regs, f)
		}
		return nil
	}
}

func WithEnableRestart() Option {
	return func(s *server) error {
		s.enableRestart = true
		return nil
	}
}

func WithDisableRestart() Option {
	return func(s *server) error {
		s.enableRestart = false
		return nil
	}
}

func WithGRPCMaxReceiveMessageSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxRecvMsgSize(size))
		}
		return nil
	}
}

func WithGRPCMaxSendMessageSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxSendMsgSize(size))
		}
		return nil
	}
}

func WithGRPCInitialWindowSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.InitialWindowSize(int32(size)))
		}
		return nil
	}
}

func WithGRPCInitialConnWindowSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.InitialConnWindowSize(int32(size)))
		}
		return nil
	}
}

func WithGRPCKeepaliveMaxConnIdle(max string) Option {
	return func(s *server) error {
		if len(max) == 0 {
			return nil
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.maxConnIdle = d
		return nil
	}
}

func WithGRPCKeepaliveMaxConnAge(max string) Option {
	return func(s *server) error {
		if len(max) == 0 {
			return nil
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.maxConnAge = d
		return nil
	}
}

func WithGRPCKeepaliveMaxConnAgeGrace(max string) Option {
	return func(s *server) error {
		if len(max) == 0 {
			return nil
		}
		d, err := timeutil.Parse(max)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.maxConnAgeGrace = d
		return nil
	}
}

func WithGRPCKeepaliveTime(dur string) Option {
	return func(s *server) error {
		if len(dur) == 0 {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.t = d
		return nil
	}
}

func WithGRPCKeepaliveTimeout(dur string) Option {
	return func(s *server) error {
		if len(dur) == 0 {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.timeout = d
		return nil
	}
}

func WithGRPCKeepaliveMinTime(min string) Option {
	return func(s *server) error {
		if len(min) == 0 {
			return nil
		}
		d, err := timeutil.Parse(min)
		if err != nil {
			return nil
		}
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.minTime = d
		return nil
	}
}

func WithGRPCKeepalivePermitWithoutStream(pws bool) Option {
	return func(s *server) error {
		if s.grpc.keepAlive == nil {
			s.grpc.keepAlive = new(grpcKeepalive)
		}
		s.grpc.keepAlive.permitWithoutStream = pws
		return nil
	}
}

func WithGRPCWriteBufferSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.WriteBufferSize(size))
		}
		return nil
	}
}

func WithGRPCReadBufferSize(size int) Option {
	return func(s *server) error {
		if size > 0 || size == -1 {
			s.grpc.opts = append(s.grpc.opts, grpc.ReadBufferSize(size))
		}
		return nil
	}
}

func WithGRPCConnectionTimeout(to string) Option {
	return func(s *server) error {
		if len(to) == 0 {
			return nil
		}
		d, err := timeutil.Parse(to)
		if err != nil {
			return nil
		}
		s.grpc.opts = append(s.grpc.opts, grpc.ConnectionTimeout(d))
		return nil
	}
}

func WithGRPCMaxHeaderListSize(size int) Option {
	return func(s *server) error {
		if size > 0 {
			s.grpc.opts = append(s.grpc.opts, grpc.MaxHeaderListSize(uint32(size)))
		}
		return nil
	}
}

func WithGRPCHeaderTableSize(size int) Option {
	return func(s *server) error {
		if size > 0 {
			s.grpc.opts = append(s.grpc.opts, grpc.HeaderTableSize(uint32(size)))
		}
		return nil
	}
}

func WithGRPCInterceptors(names ...string) Option {
	return func(s *server) error {
		for _, name := range names {
			switch strings.ToLower(name) {
			case "recoverinterceptor", "recover":
				s.grpc.opts = append(
					s.grpc.opts,
					grpc.ChainUnaryInterceptor(recover.RecoverInterceptor()),
					grpc.ChainStreamInterceptor(recover.RecoverStreamInterceptor()),
				)
			case "accessloginterceptor", "accesslog":
				s.grpc.opts = append(
					s.grpc.opts,
					grpc.ChainUnaryInterceptor(logging.AccessLogInterceptor()),
					grpc.ChainStreamInterceptor(logging.AccessLogStreamInterceptor()),
				)
			case "traceinterceptor", "trace":
				s.grpc.opts = append(
					s.grpc.opts,
					grpc.ChainUnaryInterceptor(trace.TraceInterceptor()),
					grpc.ChainStreamInterceptor(trace.TraceStreamInterceptor()),
				)
			case "metricinterceptor", "metric":
				mi, msi, err := metric.MetricInterceptors()
				if err != nil {
					return errors.NewErrCriticalOption("gRPCInterceptors", "metric", errors.Wrap(err, "failed to create Interceptor"))
				}
				s.grpc.opts = append(
					s.grpc.opts,
					grpc.ChainUnaryInterceptor(mi),
					grpc.ChainStreamInterceptor(msi),
				)
			default:
			}
		}
		return nil
	}
}
