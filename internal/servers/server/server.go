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



// Package server provides implementation of Go API for managing server flow
package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
)

type Server interface {
	Name() string
	IsRunning() bool
	ListenAndServe() <-chan error
	Shutdown(context.Context) error
}

type mode uint8

const (
	REST mode = 1 + iota
	GRPC
	GQL
)

func (m mode) String() string {
	switch m {
	case REST:
		return "REST"
	case GRPC:
		return "gRPC"
	case GQL:
		return "GraphQL"
	}
	return "unknown"
}

func Mode(m string) mode {
	switch strings.ToLower(m) {
	case "rest", "http":
		return REST
	case "grpc":
		return GRPC
	case "graphql", "gql":
		return GQL
	}
	return 0
}

type server struct {
	mode mode
	name string
	mu   sync.RWMutex
	wg   sync.WaitGroup
	http struct { // REST API
		srv     *http.Server
		h       http.Handler
		starter func(net.Listener) error
	}
	grpc struct { // gRPC API
		srv  *grpc.Server
		opts []grpc.ServerOption
		reg  func(*grpc.Server)
	}
	l            net.Listener
	tcfg         *tls.Config
	pwt          time.Duration // ProbeWaitTime
	sddur        time.Duration // Shutdown Duration
	rht          time.Duration // ReadHeaderTimeout
	rt           time.Duration // ReadTimeout
	wt           time.Duration // WriteTimeout
	it           time.Duration // IdleTimeout
	port         uint
	host         string
	running      bool
	preStartFunc func() error
	preStopFunc  func() error // PreStopFunction
}

func New(opts ...Option) (Server, error) {

	srv := new(server)

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, opt := range append(defaultOpts, opts...) {
		opt(srv)
	}

	if srv.l == nil && (srv.port != 0 || srv.host != "") {
		var err error
		srv.l, err = (&net.ListenConfig{
			Control: tcp.Control,
		}).Listen(context.Background(), "tcp",
			fmt.Sprintf("%s:%d", srv.host, srv.port))

		if err != nil {
			return nil, err
		}

		if srv.tcfg != nil {
			srv.l = tls.NewListener(srv.l, srv.tcfg)
		}
	}

	if srv.l == nil {
		return nil, errors.ErrInvalidAPIConfig
	}

	switch srv.mode {
	case REST, GQL:
		if srv.http.h == nil {
			return nil, errors.ErrInvalidAPIConfig
		}
		if srv.http.srv == nil {
			srv.http.srv = new(http.Server)
		}
		if srv.rht != 0 {
			srv.http.srv.ReadHeaderTimeout = srv.rht
		}
		if srv.rt != 0 {
			srv.http.srv.ReadTimeout = srv.rt
		}
		if srv.wt != 0 {
			srv.http.srv.WriteTimeout = srv.wt
		}
		if srv.it != 0 {
			srv.http.srv.IdleTimeout = srv.it
		}
		if srv.http.h != nil {
			srv.http.srv.Handler = srv.http.h
		}
		srv.http.starter = srv.http.srv.Serve
		if srv.tcfg != nil {
			srv.http.srv.TLSConfig = srv.tcfg
			srv.http.starter = func(l net.Listener) error {
				return srv.http.srv.ServeTLS(l, "", "")
			}
		}
		srv.http.srv.SetKeepAlivesEnabled(true)
	case GRPC:
		if srv.grpc.reg == nil {
			return nil, errors.ErrInvalidAPIConfig
		}
		if srv.grpc.srv == nil {
			srv.grpc.srv = grpc.NewServer(
				srv.grpc.opts...,
			)
		}
		srv.grpc.reg(srv.grpc.srv)
	}

	return srv, nil
}

func (s *server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *server) Name() string {
	return s.name
}

func (s *server) ListenAndServe() <-chan error {
	ech := make(chan error, 1)
	var wg sync.WaitGroup
	if !s.IsRunning() {
		s.mu.Lock()
		s.running = true
		s.mu.Unlock()
		wg.Add(1)
		s.wg.Add(1)
		errgroup.Go(safety.RecoverFunc(func() (err error) {
			defer s.wg.Done()
			defer close(ech)

			if s.preStartFunc != nil {
				log.Infof("server %s executing preStartFunc", s.name)
				err = s.preStartFunc()
				if err != nil {
					ech <- err
				}
			}

			log.Infof("%s server %s starting", s.mode.String(), s.name)
			wg.Done()
			switch s.mode {
			case REST, GQL:
				err = s.http.starter(s.l)
				if err != nil && err != http.ErrServerClosed {
					ech <- err
				}
			case GRPC:
				err = s.grpc.srv.Serve(s.l)
				if err != nil {
					ech <- err
				}
			}
			log.Infof("%s server %s stopped", s.mode.String(), s.name)
			return nil
		}))
	}
	wg.Wait()
	return ech
}

func (s *server) Shutdown(ctx context.Context) (rerr error) {
	if !s.IsRunning() {
		return nil
	}
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	log.Warnf("%s server %s shutdown process starting", s.mode.String(), s.name)
	if s.preStopFunc != nil {
		ech := make(chan error, 1)
		s.wg.Add(1)
		errgroup.Go(safety.RecoverFunc(func() (err error) {
			log.Infof("server %s executing preStopFunc", s.name)
			err = s.preStopFunc()
			ech <- err
			s.wg.Done()
			return err
		}))
		time.Sleep(s.pwt)
		err := <-ech
		close(ech)
		if err != nil {
			rerr = err
		}
	} else {
		time.Sleep(s.pwt)
	}

	log.Warnf("%s server %s is now shutting down", s.mode.String(), s.name)
	switch s.mode {
	case REST, GQL:
		sctx, scancel := context.WithTimeout(ctx, s.sddur)
		defer scancel()

		s.http.srv.SetKeepAlivesEnabled(false)

		err := s.http.srv.Shutdown(sctx)
		if err != nil && err != http.ErrServerClosed {
			rerr = errors.Wrap(rerr, err.Error())
		}

		err = sctx.Err()
		if err != nil && err != context.Canceled {
			rerr = errors.Wrap(rerr, err.Error())
		}

	case GRPC:
		s.grpc.srv.GracefulStop()

	}

	s.wg.Wait()

	return
}
