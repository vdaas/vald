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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/servers/server"

type Servers struct {
	// Server represent server configuration.
	Servers []*Server `json:"servers" yaml:"servers"`

	// HealthCheckServers represent health check server configuration
	HealthCheckServers []*Server `json:"health_check_servers" yaml:"health_check_servers"`

	// MetricsServers represent metrics exporter server such as prometheus or opentelemetly or golang's pprof server
	MetricsServers []*Server `json:"metrics_servers" yaml:"metrics_servers"`

	// StartUpStrategy represent starting order of server name
	StartUpStrategy []string `json:"startup_strategy" yaml:"startup_strategy"`
	// ShutdownStrategy represent shutdonw order of server name
	ShutdownStrategy []string `json:"shutdown_strategy" yaml:"shutdown_strategy"`

	// FullShutdownDuration represent summary duration of shutdown time
	FullShutdownDuration string `json:"full_shutdown_duration" yaml:"full_shutdown_duration"`

	// TLS represent server tls configuration.
	TLS *TLS `json:"tls" yaml:"tls"`
}

type Server struct {
	Name          string `json:"name" yaml:"name"`
	Host          string `json:"host" yaml:"host"`
	Port          uint   `json:"port" yaml:"port"`
	Mode          string `json:"mode" yaml:"mode"` // gRPC, REST, GraphQL
	ProbeWaitTime string `json:"probe_wait_time" yaml:"probe_wait_time"`
	HTTP          struct {
		ShutdownDuration  string `json:"shutdown_duration" yaml:"shutdown_duration"`
		HandlerTimeout    string `json:"handler_timeout" yaml:"handler_timeout"`
		IdleTimeout       string `json:"idle_timeout" yaml:"idle_timeout"`
		ReadHeaderTimeout string `json:"read_header_timeout" yaml:"read_header_timeout"`
		ReadTimeout       string `json:"read_timeout" yaml:"read_timeout"`
		WriteTimeout      string `json:"write_timeout" yaml:"write_timeout"`
	} `json:"http" yaml:"http"`
	GRPC struct {
		Interceptors []string `json:"interceptors" yaml:"interceptors"`
	} `json:"grpc" yaml:"grpc"`
}

func (s *Servers) Bind() *Servers {
	for i := range s.Servers {
		s.Servers[i].Bind()
	}
	for i := range s.HealthCheckServers {
		s.HealthCheckServers[i].Bind()
	}

	s.FullShutdownDuration = GetActualValue(s.FullShutdownDuration)

	for i, ss := range s.StartUpStrategy {
		s.StartUpStrategy[i] = GetActualValue(ss)
	}

	for i, ss := range s.ShutdownStrategy {
		s.ShutdownStrategy[i] = GetActualValue(ss)
	}

	if s.TLS != nil {
		s.TLS.Bind()
	}
	return s
}

func (s *Server) Bind() *Server {
	s.Name = GetActualValue(s.Name)
	s.Host = GetActualValue(s.Host)
	s.Mode = GetActualValue(s.Mode)
	s.ProbeWaitTime = GetActualValue(s.ProbeWaitTime)
	s.HTTP.HandlerTimeout = GetActualValue(s.HTTP.HandlerTimeout)
	s.HTTP.ShutdownDuration = GetActualValue(s.HTTP.ShutdownDuration)
	s.HTTP.ReadHeaderTimeout = GetActualValue(s.HTTP.ReadHeaderTimeout)
	s.HTTP.ReadTimeout = GetActualValue(s.HTTP.ReadTimeout)
	s.HTTP.WriteTimeout = GetActualValue(s.HTTP.WriteTimeout)
	s.HTTP.IdleTimeout = GetActualValue(s.HTTP.IdleTimeout)
	return s
}

func (s *Server) Opts() []server.Option {
	opts := make([]server.Option, 0, 10)
	opts = append(opts,
		server.WithName(s.Name),
		server.WithHost(s.Host),
		server.WithPort(s.Port),
		server.WithProbeWaitTime(s.ProbeWaitTime),
	)

	switch mode := server.Mode(s.Mode); mode {
	case server.REST, server.GQL:
		opts = append(opts,
			server.WithReadHeaderTimeout(s.HTTP.ReadHeaderTimeout),
			server.WithReadTimeout(s.HTTP.ReadTimeout),
			server.WithWriteTimeout(s.HTTP.WriteTimeout),
			server.WithIdleTimeout(s.HTTP.IdleTimeout),
			server.WithShutdownDuration(s.HTTP.ShutdownDuration),
			server.WithServerMode(mode),
		)
	case server.GRPC:
		opts = append(opts,
			server.WithServerMode(mode),
		)
	}

	return opts
}
