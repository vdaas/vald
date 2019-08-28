// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package servers provides implementation of Go API for managing server flow
package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
)

type Option func(*server)

var (
	defaultOpts = []Option{
		WithServerMode(REST),
	}
	HealthServerOpts = func(name, host, path string, port uint) []Option {
		return []Option{
			WithName(name),
			WithHTTPHandler(func() http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodGet {
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
