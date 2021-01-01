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

// Package config providers configuration type and load configuration logic
package config

import (
	"github.com/vdaas/vald/internal/servers/server"
)

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
	HTTP          *HTTP  `json:"http" yaml:"http"`
	GRPC          *GRPC  `json:"grpc" yaml:"grpc"`
	Restart       bool   `json:"restart" yaml:"restart"`
}

type HTTP struct {
	ShutdownDuration  string `json:"shutdown_duration" yaml:"shutdown_duration"`
	HandlerTimeout    string `json:"handler_timeout" yaml:"handler_timeout"`
	IdleTimeout       string `json:"idle_timeout" yaml:"idle_timeout"`
	ReadHeaderTimeout string `json:"read_header_timeout" yaml:"read_header_timeout"`
	ReadTimeout       string `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout      string `json:"write_timeout" yaml:"write_timeout"`
}

type GRPC struct {
	BidirectionalStreamConcurrency int            `json:"bidirectional_stream_concurrency" yaml:"bidirectional_stream_concurrency"`
	MaxReceiveMessageSize          int            `json:"max_receive_message_size" yaml:"max_receive_message_size"`
	MaxSendMessageSize             int            `json:"max_send_message_size" yaml:"max_send_message_size"`
	InitialWindowSize              int            `json:"initial_window_size" yaml:"initial_window_size"`
	InitialConnWindowSize          int            `json:"initial_conn_window_size" yaml:"initial_conn_window_size"`
	Keepalive                      *GRPCKeepalive `json:"keepalive" yaml:"keepalive"`
	WriteBufferSize                int            `json:"write_buffer_size" yaml:"write_buffer_size"`
	ReadBufferSize                 int            `json:"read_buffer_size" yaml:"read_buffer_size"`
	ConnectionTimeout              string         `json:"connection_timeout" yaml:"connection_timeout"`
	MaxHeaderListSize              int            `json:"max_header_list_size" yaml:"max_header_list_size"`
	HeaderTableSize                int            `json:"header_table_size" yaml:"header_table_size"`
	Interceptors                   []string       `json:"interceptors" yaml:"interceptors"`
}

type GRPCKeepalive struct {
	MaxConnIdle     string `json:"max_conn_idle" yaml:"max_conn_idle"`
	MaxConnAge      string `json:"max_conn_age" yaml:"max_conn_age"`
	MaxConnAgeGrace string `json:"max_conn_age_grace" yaml:"max_conn_age_grace"`
	Time            string `json:"time" yaml:"time"`
	Timeout         string `json:"timeout" yaml:"timeout"`
}

func (s *Servers) Bind() *Servers {
	check := make(map[string]struct{}, len(s.Servers)+len(s.HealthCheckServers)+len(s.MetricsServers))
	for i, srv := range s.Servers {
		if srv != nil {
			s.Servers[i].Bind()
			check[srv.Name] = struct{}{}
		}
	}

	for i, srv := range s.HealthCheckServers {
		if srv != nil {
			s.HealthCheckServers[i].Bind()
			check[srv.Name] = struct{}{}
		}
	}

	for i, srv := range s.MetricsServers {
		if srv != nil {
			s.MetricsServers[i].Bind()
			check[srv.Name] = struct{}{}
		}
	}

	s.FullShutdownDuration = GetActualValue(s.FullShutdownDuration)

	sus := make([]string, 0, len(s.StartUpStrategy))
	for _, ss := range s.StartUpStrategy {
		if _, ok := check[ss]; ok {
			sus = append(sus, GetActualValue(ss))
		}
	}
	s.StartUpStrategy = sus

	sds := make([]string, 0, len(s.ShutdownStrategy))
	for _, ss := range s.ShutdownStrategy {
		if _, ok := check[ss]; ok {
			sds = append(sds, GetActualValue(ss))
		}
	}
	s.ShutdownStrategy = sds

	if s.TLS != nil {
		s.TLS.Bind()
	} else {
		s.TLS = &TLS{
			Enabled: false,
		}
	}
	return s
}

func (s *Servers) GetGRPCStreamConcurrency() (c int) {
	for _, s := range s.Servers {
		if s.GRPC != nil {
			return s.GRPC.BidirectionalStreamConcurrency
		}
	}
	return 0
}

func (h *HTTP) Bind() *HTTP {
	h.HandlerTimeout = GetActualValue(h.HandlerTimeout)
	h.ShutdownDuration = GetActualValue(h.ShutdownDuration)
	h.ReadHeaderTimeout = GetActualValue(h.ReadHeaderTimeout)
	h.ReadTimeout = GetActualValue(h.ReadTimeout)
	h.WriteTimeout = GetActualValue(h.WriteTimeout)
	h.IdleTimeout = GetActualValue(h.IdleTimeout)
	return h
}

func (g *GRPC) Bind() *GRPC {
	g.ConnectionTimeout = GetActualValue(g.ConnectionTimeout)
	for i, ic := range g.Interceptors {
		g.Interceptors[i] = GetActualValue(ic)
	}
	return g
}

func (k *GRPCKeepalive) Bind() *GRPCKeepalive {
	k.MaxConnIdle = GetActualValue(k.MaxConnIdle)
	k.MaxConnAge = GetActualValue(k.MaxConnAge)
	k.MaxConnAgeGrace = GetActualValue(k.MaxConnAgeGrace)
	k.Time = GetActualValue(k.Time)
	k.Timeout = GetActualValue(k.Timeout)
	return k
}

func (s *Server) Bind() *Server {
	s.Name = GetActualValue(s.Name)
	s.Host = GetActualValue(s.Host)
	s.Mode = GetActualValue(s.Mode)
	s.ProbeWaitTime = GetActualValue(s.ProbeWaitTime)

	if s.HTTP != nil {
		s.HTTP.Bind()
	}

	if s.GRPC != nil {
		s.GRPC.Bind()
	}
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
			server.WithServerMode(mode),
		)
		if s.HTTP != nil {
			opts = append(opts,
				server.WithReadHeaderTimeout(s.HTTP.ReadHeaderTimeout),
				server.WithReadTimeout(s.HTTP.ReadTimeout),
				server.WithWriteTimeout(s.HTTP.WriteTimeout),
				server.WithIdleTimeout(s.HTTP.IdleTimeout),
				server.WithShutdownDuration(s.HTTP.ShutdownDuration),
			)
		}
	case server.GRPC:
		opts = append(opts,
			server.WithServerMode(mode),
		)
		if s.GRPC != nil {
			opts = append(opts,
				server.WithServerMode(mode),
				server.WithGRPCMaxReceiveMessageSize(s.GRPC.MaxReceiveMessageSize),
				server.WithGRPCMaxSendMessageSize(s.GRPC.MaxSendMessageSize),
				server.WithGRPCInitialWindowSize(s.GRPC.InitialWindowSize),
				server.WithGRPCInitialConnWindowSize(s.GRPC.InitialConnWindowSize),
				server.WithGRPCWriteBufferSize(s.GRPC.WriteBufferSize),
				server.WithGRPCReadBufferSize(s.GRPC.ReadBufferSize),
				server.WithGRPCConnectionTimeout(s.GRPC.ConnectionTimeout),
				server.WithGRPCMaxHeaderListSize(s.GRPC.MaxHeaderListSize),
				server.WithGRPCHeaderTableSize(s.GRPC.HeaderTableSize),
				server.WithGRPCInterceptors(s.GRPC.Interceptors...),
			)
			if s.GRPC.Keepalive != nil {
				opts = append(opts,
					server.WithGRPCKeepaliveMaxConnIdle(s.GRPC.Keepalive.MaxConnIdle),
					server.WithGRPCKeepaliveMaxConnAge(s.GRPC.Keepalive.MaxConnAge),
					server.WithGRPCKeepaliveMaxConnAgeGrace(s.GRPC.Keepalive.MaxConnAgeGrace),
					server.WithGRPCKeepaliveTime(s.GRPC.Keepalive.Time),
					server.WithGRPCKeepaliveTimeout(s.GRPC.Keepalive.Timeout),
				)
			}
		}
	}

	return opts
}
