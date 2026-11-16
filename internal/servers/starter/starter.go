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

package starter

import (
	"fmt"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/http/metrics"
	"github.com/vdaas/vald/internal/servers"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/tls"
)

type Server servers.Listener

type srvs struct {
	rest    func(cfg *config.Server) []server.Option
	gql     func(cfg *config.Server) []server.Option
	grpc    func(cfg *config.Server) []server.Option
	cfg     *config.Servers
	pstartf map[string]func() error
	pstopf  map[string]func() error
}

func New(sopts ...Option) (Server, error) {
	ss := &srvs{
		cfg:     new(config.Servers),
		pstartf: make(map[string]func() error, len(sopts)),
		pstopf:  make(map[string]func() error, len(sopts)),
	}

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

	if ss.cfg.TLS != nil && ss.cfg.TLS.Enabled {
		var err error
		cfg, err = tls.NewServerConfig(ss.cfg.TLS.Opts()...)
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
			restOpts := sc.Opts()
			if s.rest != nil {
				restOpts = append(restOpts, s.rest(sc)...)
			}
			srv, err := server.New(
				append(append(restOpts,
					server.WithTLSConfig(cfg),
				), s.hookOpts(sc.Name)...)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		case server.GRPC:
			grpcOpts := append(
				sc.Opts(),
				server.WithTLSConfig(cfg),
			)

			if s.grpc != nil {
				grpcOpts = append(grpcOpts, s.grpc(sc)...)
			}
			grpcOpts = append(grpcOpts, s.hookOpts(sc.Name)...)

			srv, err := server.New(grpcOpts...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		case server.GQL:
			gqlOpts := sc.Opts()
			if s.gql != nil {
				gqlOpts = append(gqlOpts, s.gql(sc)...)
			}
			srv, err := server.New(
				append(append(gqlOpts,
					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				), s.hookOpts(sc.Name)...)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		}
	}

	return opts, nil
}

// hookOpts returns the pre-start and pre-stop functions registered for the server name
// by WithPreStartFunc and WithPreStopFunc. They are appended last so that they override any default hook.
func (s *srvs) hookOpts(name string) (opts []server.Option) {
	if f, ok := s.pstartf[name]; ok {
		opts = append(opts, server.WithPreStartFunc(f))
	}
	if f, ok := s.pstopf[name]; ok {
		opts = append(opts, server.WithPreStopFunction(f))
	}
	return opts
}

func (s *srvs) setupHealthCheck(*tls.Config) ([]servers.Option, error) {
	opts := make([]servers.Option, 0, len(s.cfg.HealthCheckServers))
	for _, hsc := range s.cfg.HealthCheckServers {
		srv, err := server.New(append(healthServerOpts(hsc), s.hookOpts(hsc.Name)...)...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, servers.WithServer(srv))
	}
	return opts, nil
}

// healthServerOpts returns the HealthServerOpts presets overridden by hsc.
// hsc.Opts passes an unset ("") duration through, which the options parse as 0
// and so wipe the presets (the default Helm values leave the health servers'
// http timeouts empty). The presets are therefore applied again after
// hsc.Opts, followed by the config values they overrode: network and mode, and
// only the durations the config actually sets.
func healthServerOpts(hsc *config.Server) []server.Option {
	preset := func() []server.Option {
		return server.HealthServerOpts(
			hsc.Name,
			hsc.Host,
			fmt.Sprintf("/%s", strings.ToLower(hsc.Name)),
			hsc.Port)
	}
	opts := append(preset(), hsc.Opts()...)
	opts = append(opts, preset()...)
	mode := server.Mode(hsc.Mode)
	opts = append(opts,
		server.WithNetwork(hsc.Network),
		server.WithServerMode(mode),
	)
	type durOpt struct {
		opt func(string) server.Option
		dur string
	}
	durs := []durOpt{{server.WithProbeWaitTime, hsc.ProbeWaitTime}}
	// hsc.Opts applies the http durations only in these modes.
	if h := hsc.HTTP; h != nil && (mode == server.REST || mode == server.GQL) {
		durs = append(durs,
			durOpt{server.WithReadHeaderTimeout, h.ReadHeaderTimeout},
			durOpt{server.WithReadTimeout, h.ReadTimeout},
			durOpt{server.WithWriteTimeout, h.WriteTimeout},
			durOpt{server.WithIdleTimeout, h.IdleTimeout},
			durOpt{server.WithShutdownDuration, h.ShutdownDuration},
		)
	}
	for _, d := range durs {
		if d.dur != "" {
			opts = append(opts, d.opt(d.dur))
		}
	}
	return opts
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
				append(append(msc.Opts(),
					hopt,
					server.WithTLSConfig(cfg),
					server.WithPreStartFunc(func() error {
						return nil
					}),
					server.WithPreStopFunction(func() error {
						return nil
					}),
				), s.hookOpts(msc.Name)...)...)
			if err != nil {
				return nil, err
			}
			opts = append(opts, servers.WithServer(srv))
		}
	}
	return opts, nil
}
