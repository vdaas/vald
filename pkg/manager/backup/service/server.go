//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package service manages the main logic of server.
package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vdaas/vald/apis/grpc/manager/backup"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/http/metrics"
	"github.com/vdaas/vald/internal/servers"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/tls"
	"google.golang.org/grpc"
)

type Server servers.Listener

type srvs struct {
	rest http.Handler
	gql  http.Handler
	grpc backup.BackupServer
	cfg  *config.Servers
}

func NewServer(sopts ...Option) (Server, error) {
	ss := new(srvs)
	for _, opt := range sopts {
		opt(ss)
	}

	opts := make([]servers.Option, 0, 3+
		len(ss.cfg.Servers)+
		len(ss.cfg.HealthCheckServers)+
		len(ss.cfg.MetricsServers))

	opts = append(opts,
		servers.WithShutdownDuration(ss.cfg.FullShutdownDuration),
		servers.WithStartUpStrategy(ss.cfg.StartUpStrategy),
		servers.WithShutdownStrategy(ss.cfg.ShutdownStrategy))

	var cfg *tls.Config

	if ss.cfg.TLS.Enabled {
		var err error
		cfg, err = tls.New(
			tls.WithCert(ss.cfg.TLS.Cert),
			tls.WithKey(ss.cfg.TLS.Key),
			tls.WithCa(ss.cfg.TLS.CA),
		)
		if err != nil {
			return nil, err
		}
	}

	apiOpts, err := ss.setupAPIs(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, apiOpts...)

	hcOpts, err := ss.setupHealthCheck(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, hcOpts...)

	mOpts, err := ss.setupMetrics(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, mOpts...)

	return servers.New(opts...), nil
}

func (s *srvs) setupAPIs(cfg *tls.Config) ([]servers.Option, error) {
	opts := make([]servers.Option, 0, len(s.cfg.Servers))
	for _, sc := range s.cfg.Servers {
		switch mode := server.Mode(sc.Mode); mode {
		case server.REST:
			srv, err := server.New(
				append(sc.Opts(),
					server.WithHTTPHandler(s.rest),
					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		case server.GRPC:
			gopts := make([]grpc.ServerOption, 0, len(sc.GRPC.Interceptors))
			for _, ic := range sc.GRPC.Interceptors {
				switch strings.ToLower(ic) {
				case "valid", "validate", "validation":
					// TODO create interceptor in internal
					// TODO add grpc interceptor in internal
				}
			}
			srv, err := server.New(
				append(sc.Opts(),
					server.WithGRPCRegistFunc(func(gsrv *grpc.Server) {
						backup.RegisterBackupServer(gsrv, s.grpc)
					}),
					server.WithGRPCOption(gopts[:len(gopts)]...),

					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		case server.GQL:
			srv, err := server.New(
				append(sc.Opts(),
					server.WithHTTPHandler(s.gql),
					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		}
	}

	return opts, nil
}

func (s *srvs) setupHealthCheck(cfg *tls.Config) ([]servers.Option, error) {
	opts := make([]servers.Option, 0, len(s.cfg.HealthCheckServers))
	for _, hsc := range s.cfg.HealthCheckServers {
		srv, err := server.New(
			append(server.HealthServerOpts(
				hsc.Name,
				hsc.Host,
				fmt.Sprintf("/%s", strings.ToLower(hsc.Name)),
				hsc.Port),
				hsc.Opts()...)...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, servers.WithServer(srv))
	}
	return opts, nil
}

func (s *srvs) setupMetrics(cfg *tls.Config) ([]servers.Option, error) {
	opts := make([]servers.Option, 0, len(s.cfg.MetricsServers))
	for _, msc := range s.cfg.MetricsServers {
		var hopt server.Option
		switch strings.ToLower(msc.Name) {
		case "prof", "pprof", "profile", "profiler":
			hopt = server.WithHTTPHandler(metrics.NewPProfHandler())
		default:
			continue
		}
		if hopt != nil {
			srv, err := server.New(
				append(msc.Opts(),
					hopt,
					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		}
	}
	return opts, nil
}
