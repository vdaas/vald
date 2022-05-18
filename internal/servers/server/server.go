//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package server provides implementation of Go API for managing server flow
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/credentials"
	"github.com/vdaas/vald/internal/net/grpc/keepalive"
	glog "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/net/quic"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
)

type Server interface {
	Name() string
	IsRunning() bool
	ListenAndServe(context.Context, chan<- error) error
	Shutdown(context.Context) error
}

type ServerMode uint8

const (
	REST ServerMode = 1 + iota
	GRPC
	GQL
)

func (m ServerMode) String() string {
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

func Mode(m string) ServerMode {
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
	mode ServerMode
	name string
	mu   sync.RWMutex
	wg   sync.WaitGroup
	eg   errgroup.Group
	http struct { // REST API
		srv     *http.Server
		h       http.Handler
		starter func(net.Listener) error
	}
	grpc struct { // gRPC API
		srv       *grpc.Server
		keepAlive *grpcKeepalive
		opts      []grpc.ServerOption
		regs      []func(*grpc.Server)
	}
	lc            *net.ListenConfig
	tcfg          *tls.Config
	pwt           time.Duration // ProbeWaitTime
	sddur         time.Duration // Shutdown Duration
	rht           time.Duration // ReadHeaderTimeout
	rt            time.Duration // ReadTimeout
	wt            time.Duration // WriteTimeout
	it            time.Duration // IdleTimeout
	ctrl          control.SocketController
	sockFlg       control.SocketFlag
	network       net.NetworkType
	socketPath    string
	port          uint16
	host          string
	enableRestart bool
	shuttingDown  bool
	running       bool
	preStartFunc  func() error
	preStopFunc   func() error // PreStopFunction
}

type grpcKeepalive struct {
	maxConnIdle         time.Duration
	maxConnAge          time.Duration
	maxConnAgeGrace     time.Duration
	t                   time.Duration
	timeout             time.Duration
	minTime             time.Duration
	permitWithoutStream bool
}

func New(opts ...Option) (Server, error) {
	srv := new(server)

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, opt := range append(defaultOptions, opts...) {
		opt(srv)
	}
	if srv.eg == nil {
		log.Warnf("errgroup not found for %s, getting new errgroup.", srv.name)
		srv.eg = errgroup.Get()
	}

	var keepAlive time.Duration
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
		if srv.tcfg != nil &&
			(len(srv.tcfg.Certificates) != 0 ||
				srv.tcfg.GetCertificate != nil ||
				srv.tcfg.GetConfigForClient != nil) {
			srv.http.srv.TLSConfig = srv.tcfg
			srv.http.starter = func(l net.Listener) error {
				return srv.http.srv.ServeTLS(l, "", "")
			}
		}
		srv.http.srv.SetKeepAlivesEnabled(true)
	case GRPC:
		if srv.grpc.regs == nil {
			return nil, errors.ErrInvalidAPIConfig
		}

		if srv.grpc.keepAlive != nil {
			srv.grpc.opts = append(srv.grpc.opts,
				grpc.KeepaliveParams(keepalive.ServerParameters{
					MaxConnectionIdle:     srv.grpc.keepAlive.maxConnIdle,
					MaxConnectionAge:      srv.grpc.keepAlive.maxConnAge,
					MaxConnectionAgeGrace: srv.grpc.keepAlive.maxConnAgeGrace,
					Time:                  srv.grpc.keepAlive.t,
					Timeout:               srv.grpc.keepAlive.timeout,
				}),
				grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
					MinTime:             srv.grpc.keepAlive.minTime,
					PermitWithoutStream: srv.grpc.keepAlive.permitWithoutStream,
				}),
			)
			keepAlive = srv.grpc.keepAlive.t
		}

		if srv.tcfg != nil &&
			(len(srv.tcfg.Certificates) != 0 ||
				srv.tcfg.GetCertificate != nil ||
				srv.tcfg.GetConfigForClient != nil) {
			srv.grpc.opts = append(srv.grpc.opts,
				grpc.Creds(credentials.NewTLS(srv.tcfg)),
			)
		}

		if srv.grpc.srv == nil {
			srv.grpc.srv = grpc.NewServer(
				srv.grpc.opts...,
			)
		}
		for _, reg := range srv.grpc.regs {
			reg(srv.grpc.srv)
		}
	}

	if srv.lc == nil {
		srv.ctrl = control.New(srv.sockFlg, int(keepAlive))
		srv.lc = &net.ListenConfig{
			KeepAlive: keepAlive,
			Control: func(network, addr string, c syscall.RawConn) (err error) {
				if srv.ctrl != nil {
					return srv.ctrl.GetControl()(network, addr, c)
				}
				log.Warn("socket controller is nil")
				return nil
			},
		}
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

func (s *server) ListenAndServe(ctx context.Context, ech chan<- error) (err error) {
	if !s.IsRunning() {
		s.mu.Lock()
		s.running = true
		s.mu.Unlock()

		if s.preStartFunc != nil {
			log.Infof("server %s executing preStartFunc", s.name)
			err = s.preStartFunc()
			if err != nil {
				return err
			}
		}
		network := func() string {
			if s.network == 0 || s.network == net.Unknown || strings.EqualFold(s.network.String(), net.Unknown.String()) {
				return net.TCP.String()
			}
			return s.network.String()
		}()

		addr := func() string {
			if s.network == net.UNIX {
				if len(s.socketPath) == 0 {
					s.socketPath = os.TempDir() + string(os.PathSeparator) + s.name + "." + strconv.Itoa(os.Getpid()) + ".sock"
				}
				return s.socketPath
			}
			return net.JoinHostPort(s.host, s.port)
		}()

		var l net.Listener
		if net.IsUDP(network) {
			l, err = quic.Listen(ctx, addr, s.tcfg)
			if err != nil {
				log.Errorf("failed to listen udp socket for quic:\terror %v ", err)
				return err
			}
		} else {
			l, err = s.lc.Listen(ctx, network, addr)
			if err != nil {
				log.Errorf("failed to listen socket %v", err)
				return err
			}
			var file *os.File
			switch lt := l.(type) {
			case *net.TCPListener:
				file, err = lt.File()
				if err != nil {
					log.Errorf("failed to listen tcp socket %v", err)
					return err
				}
			case *net.UnixListener:
				file, err = lt.File()
				if err != nil {
					log.Errorf("failed to listen unix socket %v", err)
					return err
				}
			}
			if file != nil {
				err = syscall.SetNonblock(int(file.Fd()), true)
				if err != nil {
					return err
				}
			}
			if s.tcfg != nil &&
				(len(s.tcfg.Certificates) != 0 ||
					s.tcfg.GetCertificate != nil ||
					s.tcfg.GetConfigForClient != nil) {
				l = tls.NewListener(l, s.tcfg)
			}
		}

		if l == nil {
			return errors.ErrInvalidAPIConfig
		}

		s.wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer s.wg.Done()
			for {
				if !s.IsRunning() {
					s.mu.Lock()
					s.running = true
					s.mu.Unlock()
				}
				log.Infof("%s server %s starting on %s://%s", s.mode.String(), s.name, l.Addr().Network(), l.Addr().String())

				switch s.mode {
				case REST, GQL:
					err = s.http.starter(l)
					if err != nil &&
						!errors.Is(err, http.ErrServerClosed) &&
						!errors.Is(err, context.Canceled) &&
						!errors.Is(err, context.DeadlineExceeded) {
						select {
						case <-ctx.Done():
							log.Error(errors.Wrap(ctx.Err(), err.Error()))
						case ech <- err:
						}
					}
				case GRPC:
					glog.Init()
					err = s.grpc.srv.Serve(l)
					if err != nil &&
						!errors.Is(err, grpc.ErrServerStopped) &&
						!errors.Is(err, context.Canceled) &&
						!errors.Is(err, context.DeadlineExceeded) {
						select {
						case <-ctx.Done():
							log.Error(errors.Wrap(ctx.Err(), err.Error()))
						case ech <- err:
						}
					}
				}
				err = nil
				s.mu.Lock()
				s.running = false
				s.mu.Unlock()

				s.mu.RLock()
				if !s.enableRestart || s.shuttingDown {
					s.mu.RUnlock()
					return
				}
				s.mu.RUnlock()
				log.Infof("%s server %s stopped", s.mode.String(), s.name)
			}
		}))
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) (rerr error) {
	if !s.IsRunning() {
		return nil
	}
	s.mu.Lock()
	s.running = false
	s.shuttingDown = true
	s.mu.Unlock()

	log.Warnf("%s server %s shutdown process starting", s.mode.String(), s.name)
	if s.preStopFunc != nil {
		ech := make(chan error, 1)
		s.wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer close(ech)
			log.Infof("server %s executing preStopFunc", s.name)
			err = s.preStopFunc()
			if err != nil {
				select {
				case <-ctx.Done():
				case ech <- nil:
				}
			}
			s.wg.Done()
			select {
			case <-ctx.Done():
			case ech <- nil:
			}
			return nil
		}))
		select {
		case <-ctx.Done():
		case <-time.After(s.pwt):
		case err := <-ech:
			if err != nil {
				rerr = err
			}
		}

	} else {
		select {
		case <-ctx.Done():
		case <-time.After(s.pwt):
		}
	}

	if len(s.socketPath) != 0 {
		defer func() {
			err := os.RemoveAll(s.socketPath)
			if err != nil {
				rerr = errors.Wrap(rerr, err.Error())
			}
		}()
	}

	log.Warnf("%s server %s is now shutting down", s.mode.String(), s.name)
	switch s.mode {
	case REST, GQL:
		sctx, scancel := context.WithTimeout(ctx, s.sddur)
		defer scancel()
		s.http.srv.SetKeepAlivesEnabled(false)
		err := s.http.srv.Shutdown(sctx)
		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) &&
			!errors.Is(err, grpc.ErrServerStopped) &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			rerr = errors.Wrap(rerr, err.Error())
		}
		err = sctx.Err()
		if err != nil &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) {
			rerr = errors.Wrap(rerr, err.Error())
		}

	case GRPC:
		s.grpc.srv.GracefulStop()
	}

	s.wg.Wait()

	return
}
