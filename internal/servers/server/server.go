// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/control"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/credentials"
	"github.com/vdaas/vald/internal/net/grpc/health"
	"github.com/vdaas/vald/internal/net/grpc/keepalive"
	_ "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/net/grpc/stats"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/tls"
	"golang.org/x/net/http2"
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
	http struct {
		h        http.Handler
		srv      *http.Server
		h2srv    *http2.Server
		enableH2 bool
		starter  func(net.Listener) error
	}
	ctrl         control.SocketController
	eg           errgroup.Group
	preStopFunc  func() error
	preStartFunc func() error
	lc           *net.ListenConfig
	tcfg         *tls.Config
	name         string
	host         string
	socketPath   string
	listener     net.Listener
	grpc         struct {
		srv        *grpc.Server
		keepAlive  *grpcKeepalive
		maxMsgSize int
		opts       []grpc.ServerOption
		regs       []func(*grpc.Server)
	}
	wg            sync.WaitGroup
	network       net.NetworkType
	rt            time.Duration
	it            time.Duration
	sddur         time.Duration
	sockFlg       control.SocketFlag
	rht           time.Duration
	pwt           time.Duration
	wt            time.Duration
	mu            sync.RWMutex
	port          uint16
	shuttingDown  bool
	running       bool
	enableRestart bool
	mode          ServerMode
}

// restartInterval is the wait before a stopped server is restarted, which keeps a persistently failing server from busy looping.
const restartInterval = time.Second

type grpcKeepalive struct {
	maxConnIdle         time.Duration
	maxConnAge          time.Duration
	maxConnAgeGrace     time.Duration
	t                   time.Duration
	timeout             time.Duration
	minTime             time.Duration
	permitWithoutStream bool
}

// New returns Server implementation.
// New creates and returns a configured Server using the supplied Option functions.
// It applies options (defaults first), initializes missing internal components (for example the error group),
// and configures mode-specific internals for REST/GraphQL or gRPC.
// It validates required configuration and returns an error if an option fails or required handlers/registrations are missing.
func New(opts ...Option) (Server, error) {
	srv := new(server)

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(srv); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
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
		srv.http.srv.Handler = srv.http.h
		srv.http.starter = srv.http.srv.Serve
		srv.http.srv.SetKeepAlivesEnabled(true)
		if srv.tcfg != nil &&
			(len(srv.tcfg.Certificates) != 0 ||
				srv.tcfg.GetCertificate != nil ||
				srv.tcfg.GetConfigForClient != nil) {
			srv.http.srv.TLSConfig = srv.tcfg
			srv.http.starter = func(l net.Listener) (err error) {
				if srv.http.enableH2 {
					if srv.http.h2srv != nil {
						srv.http.h2srv.IdleTimeout = srv.it
					}
					err = http2.ConfigureServer(srv.http.srv, srv.http.h2srv)
					if err != nil {
						return err
					}
				}
				err = srv.http.srv.ServeTLS(l, "", "")
				if err != nil {
					return err
				}
				return nil
			}
		}
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

		if srv.grpc.srv != nil {
			srv.grpc.srv.GracefulStop()
		}
		grpc.InitCodec(srv.grpc.maxMsgSize)
		srv.grpc.srv = grpc.NewServer(
			srv.grpc.opts...,
		)
		for _, reg := range srv.grpc.regs {
			reg(srv.grpc.srv)
		}
		stats.Register(srv.grpc.srv)
		health.Register(srv.grpc.srv)
	}
	if srv.lc == nil {
		srv.ctrl = control.New(srv.sockFlg, control.KeepAliveSeconds(keepAlive))
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

// skipcq: GO-R1005
func (s *server) ListenAndServe(ctx context.Context, ech chan<- error) (err error) {
	if !s.IsRunning() {
		s.mu.Lock()
		s.running = true
		s.mu.Unlock()
		// reset the running flag when the server failed to start so that it is neither reported as running nor skipped on retry.
		defer func() {
			if err != nil {
				s.mu.Lock()
				s.running = false
				s.mu.Unlock()
			}
		}()

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
				if s.socketPath == "" {
					sockFile := strings.Join([]string{s.name, strconv.Itoa(os.Getpid()), "sock"}, ".")
					s.socketPath = file.Join(os.TempDir(), sockFile)
				}
				return s.socketPath
			}
			return net.JoinHostPort(s.host, s.port)
		}()
		var l net.Listener
		l, err = s.listen(ctx, network, addr)
		if err != nil {
			return err
		}

		s.mu.Lock()
		s.listener = l
		s.mu.Unlock()

		s.wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer s.wg.Done()
			// running stays set while the server is restarted so that a Shutdown issued in between is not skipped.
			defer func() {
				s.mu.Lock()
				s.running = false
				s.mu.Unlock()
			}()
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
					if errors.Is(err, http.ErrServerClosed) {
						err = nil
					}
				case GRPC:
					err = s.grpc.srv.Serve(l)
					if errors.Is(err, grpc.ErrServerStopped) {
						err = nil
					}
				}
				if err != nil {
					s.mu.RLock()
					restart := s.enableRestart
					s.mu.RUnlock()
					if restart {
						// the server is restarted below; repeated sends would fill ech.
						log.Warnf("%s server %s stopped with error: %v", s.mode.String(), s.name, err)
					} else {
						// bounded by ctx so that the goroutine does not block forever once
						// nobody reads ech; ech must stay open until this goroutine exits.
						select {
						case <-ctx.Done():
						case ech <- err:
						}
					}
				}
				err = nil

				// Serve closes l when it returns, so a restart needs a new listener.
				for {
					s.mu.RLock()
					stop := !s.enableRestart || s.shuttingDown
					s.mu.RUnlock()
					if stop {
						return nil
					}
					log.Infof("%s server %s stopped, restarting in %s", s.mode.String(), s.name, restartInterval.String())
					timer := time.NewTimer(restartInterval)
					select {
					case <-ctx.Done():
						timer.Stop()
						return nil
					case <-timer.C:
					}
					// hold the lock while re-listening so that Shutdown either sees the new listener or stops the restart.
					s.mu.Lock()
					if s.shuttingDown {
						s.mu.Unlock()
						return nil
					}
					l, err = s.listen(ctx, network, addr)
					if err == nil {
						s.listener = l
					}
					s.mu.Unlock()
					if err == nil {
						break
					}
					// not sent to ech: a restarting server only logs its errors.
					log.Warnf("%s server %s failed to listen for restart: %v", s.mode.String(), s.name, err)
					err = nil
				}
			}
		}))
	}
	return nil
}

// listen opens a plain listener on addr; TLS is terminated by the HTTP or gRPC server itself.
func (s *server) listen(ctx context.Context, network, addr string) (l net.Listener, err error) {
	if s.tcfg != nil && net.IsUDP(network) {
		log.Error("QUIC protocol is not supported yet")
		return nil, errors.ErrUnsupportedClientMethod
	}
	if net.IsUDP(network) {
		network = net.TCP.String()
	}
	l, err = s.lc.Listen(ctx, network, addr)
	if err != nil {
		log.Errorf("failed to listen socket %v", err)
		return nil, err
	}
	// The listener is already non-blocking and registered with the netpoller. It is
	// deliberately not dup'ed via File(): (*os.File).Fd on that dup switches the shared
	// open file description, and so this listener, to blocking mode.
	// TLS is not applied here: ServeTLS (REST/GQL) and grpc.Creds (GRPC) set up in New already
	// terminate TLS, and wrapping the listener as well would require clients to handshake twice.
	return l, nil
}

// skipcq: GO-R1005
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
			defer s.wg.Done()
			defer close(ech)
			log.Infof("server %s executing preStopFunc", s.name)
			err = s.preStopFunc()
			select {
			case <-ctx.Done():
			case ech <- err:
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

	if s.socketPath != "" {
		defer func() {
			err := os.RemoveAll(s.socketPath)
			if err != nil {
				rerr = errors.Join(rerr, err)
			}
		}()
	}

	log.Warnf("%s server %s is now shutting down", s.mode.String(), s.name)
	switch s.mode {
	case REST, GQL:
		// a zero shutdown duration means "no per-server bound": an already
		// expired context would skip draining in-flight requests entirely,
		// so rely on the caller's ctx alone in that case.
		sctx := ctx
		if s.sddur > 0 {
			var scancel context.CancelFunc
			sctx, scancel = context.WithTimeout(ctx, s.sddur)
			defer scancel()
		}
		s.http.srv.SetKeepAlivesEnabled(false)
		err := s.http.srv.Shutdown(sctx)
		if errors.IsNot(err, http.ErrServerClosed, grpc.ErrServerStopped, context.Canceled, context.DeadlineExceeded) {
			rerr = errors.Join(rerr, err)
		}

	case GRPC:
		// GracefulStop blocks until every pending RPC (including long-lived
		// streams) finishes, so bound it by ctx and force Stop on expiry.
		// A plain goroutine, not s.eg: a limited errgroup may run f inline or drop
		// it once its context is cancelled, which would leave done open forever.
		done := make(chan struct{})
		go func() {
			defer close(done)
			s.grpc.srv.GracefulStop()
		}()
		select {
		case <-done:
		case <-ctx.Done():
			log.Warnf("%s server %s graceful stop did not finish before the shutdown deadline, forcing stop: %v", s.mode.String(), s.name, context.Cause(ctx))
			// Stop closes every connection, which also unblocks GracefulStop.
			s.grpc.srv.Stop()
			<-done
		}
	}

	s.mu.Lock()
	if s.listener != nil {
		err := s.listener.Close()
		// the underlying server (http.Server / grpc.Server) has already closed the listener, so a second Close is not an error.
		if err != nil && !errors.Is(err, os.ErrClosed) && !errors.Is(err, net.ErrClosed) {
			rerr = errors.Join(rerr, err)
		}
		s.listener = nil
	}
	s.mu.Unlock()

	s.wg.Wait()

	return rerr
}
