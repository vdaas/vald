//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package config

import (
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/admin"
	"github.com/vdaas/vald/internal/net/grpc/reflection"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/strings"
)

// Servers represents the configuration of server list.
type Servers struct {
	// TLS represent server tls configuration.
	TLS *TLS `json:"tls" yaml:"tls"`

	// FullShutdownDuration represent summary duration of shutdown time
	FullShutdownDuration string `json:"full_shutdown_duration" yaml:"full_shutdown_duration"`

	// Server represent server configuration.
	Servers []*Server `json:"servers" yaml:"servers"`

	// HealthCheckServers represent health check server configuration
	HealthCheckServers []*Server `json:"health_check_servers" yaml:"health_check_servers"`

	// MetricsServers represent metrics exporter server such as golang's pprof server
	MetricsServers []*Server `json:"metrics_servers" yaml:"metrics_servers"`

	// StartUpStrategy represent starting order of server name
	StartUpStrategy []string `json:"startup_strategy" yaml:"startup_strategy"`

	// ShutdownStrategy represent shutdown order of server name
	ShutdownStrategy []string `json:"shutdown_strategy" yaml:"shutdown_strategy"`
}

// Server represents the server configuration.
type Server struct {
	GRPC          *GRPC         `json:"grpc,omitempty"            yaml:"grpc"`
	SocketOption  *SocketOption `json:"socket_option,omitempty"   yaml:"socket_option"`
	HTTP          *HTTP         `json:"http,omitempty"            yaml:"http"`
	Name          string        `json:"name,omitempty"            yaml:"name"`
	Network       string        `json:"network,omitempty"         yaml:"network"`
	Host          string        `json:"host,omitempty"            yaml:"host"`
	SocketPath    string        `json:"socket_path,omitempty"     yaml:"socket_path"`
	Mode          string        `json:"mode,omitempty"            yaml:"mode"` // gRPC, REST, GraphQL
	ProbeWaitTime string        `json:"probe_wait_time,omitempty" yaml:"probe_wait_time"`
	Restart       bool          `json:"restart,omitempty"         yaml:"restart"`
	Port          uint16        `json:"port,omitempty"            yaml:"port"`
}

// HTTP represents the configuration for HTTP.
type HTTP struct {
	HTTP2 *HTTP2 `json:"http2"               yaml:"http2"`
	// ShutdownDuration represents the duration for the http server to shutdown.
	ShutdownDuration string `json:"shutdown_duration"   yaml:"shutdown_duration"`
	// HandlerTimeout represents the timeout duration for http handlers.
	HandlerTimeout    string `json:"handler_timeout"     yaml:"handler_timeout"`
	IdleTimeout       string `json:"idle_timeout"        yaml:"idle_timeout"`
	ReadHeaderTimeout string `json:"read_header_timeout" yaml:"read_header_timeout"`
	ReadTimeout       string `json:"read_timeout"        yaml:"read_timeout"`
	WriteTimeout      string `json:"write_timeout"       yaml:"write_timeout"`
}

// HTTP2 represents the configuration for HTTP2.
type HTTP2 struct {
	HandlerLimit                 int    `json:"handler_limit,omitempty"                    yaml:"handler_limit"`
	Enabled                      bool   `json:"enabled,omitempty"                          yaml:"enabled"`
	PermitProhibitedCipherSuites bool   `json:"permit_prohibited_cipher_suites,omitempty"  yaml:"permit_prohibited_cipher_suites"`
	MaxUploadBufferPerConnection int32  `json:"max_upload_buffer_per_connection,omitempty" yaml:"max_upload_buffer_per_connection"`
	MaxUploadBufferPerStream     int32  `json:"max_upload_buffer_per_stream,omitempty"     yaml:"max_upload_buffer_per_stream"`
	MaxConcurrentStreams         uint32 `json:"max_concurrent_streams,omitempty"           yaml:"max_concurrent_streams"`
	MaxDecoderHeaderTableSize    uint32 `json:"max_decoder_header_table_size,omitempty"    yaml:"max_decoder_header_table_size"`
	MaxEncoderHeaderTableSize    uint32 `json:"max_encoder_header_table_size,omitempty"    yaml:"max_encoder_header_table_size"`
	MaxReadFrameSize             uint32 `json:"max_read_frame_size,omitempty"              yaml:"max_read_frame_size"`
}

// Bind binds the actual data from the HTTP2 receiver fields.
func (h *HTTP2) Bind() *HTTP2 {
	// No fields to bind as per rules
	return h
}

// GRPC represents the configuration for gPRC.
type GRPC struct {
	EnableAdmin                    bool           `json:"enable_admin,omitempty"                     yaml:"enable_admin"`
	EnableChannelz                 bool           `json:"enable_channelz,omitempty"                  yaml:"enable_channelz"`
	EnableReflection               bool           `json:"enable_reflection,omitempty"                yaml:"enable_reflection"`
	SharedWriteBuffer              bool           `json:"shared_write_buffer,omitempty"              yaml:"shared_write_buffer"`
	WaitForHandlers                bool           `json:"wait_for_handlers,omitempty"                yaml:"wait_for_handlers"`
	HeaderTableSize                uint32         `json:"header_table_size,omitempty"                yaml:"header_table_size"`
	MaxConcurrentStreams           uint32         `json:"max_concurrent_streams,omitempty"           yaml:"max_concurrent_streams"`
	MaxHeaderListSize              uint32         `json:"max_header_list_size,omitempty"             yaml:"max_header_list_size"`
	NumStreamWorkers               uint32         `json:"num_stream_workers,omitempty"               yaml:"num_stream_workers"`
	BidirectionalStreamConcurrency int            `json:"bidirectional_stream_concurrency,omitempty" yaml:"bidirectional_stream_concurrency"`
	InitialConnWindowSize          int            `json:"initial_conn_window_size,omitempty"         yaml:"initial_conn_window_size"`
	InitialWindowSize              int            `json:"initial_window_size,omitempty"              yaml:"initial_window_size"`
	MaxReceiveMessageSize          int            `json:"max_receive_message_size,omitempty"         yaml:"max_receive_message_size"`
	MaxSendMessageSize             int            `json:"max_send_message_size,omitempty"            yaml:"max_send_message_size"`
	ReadBufferSize                 int            `json:"read_buffer_size,omitempty"                 yaml:"read_buffer_size"`
	WriteBufferSize                int            `json:"write_buffer_size,omitempty"                yaml:"write_buffer_size"`
	ConnectionTimeout              string         `json:"connection_timeout,omitempty"               yaml:"connection_timeout"`
	Interceptors                   []string       `json:"interceptors,omitempty"                     yaml:"interceptors"`
	Keepalive                      *GRPCKeepalive `json:"keepalive,omitempty"                        yaml:"keepalive"`
}

// GRPCKeepalive represents the configuration for gRPC keep-alive.
type GRPCKeepalive struct {
	MaxConnIdle         string `json:"max_conn_idle"         yaml:"max_conn_idle"`
	MaxConnAge          string `json:"max_conn_age"          yaml:"max_conn_age"`
	MaxConnAgeGrace     string `json:"max_conn_age_grace"    yaml:"max_conn_age_grace"`
	Time                string `json:"time"                  yaml:"time"`
	Timeout             string `json:"timeout"               yaml:"timeout"`
	MinTime             string `json:"min_time"              yaml:"min_time"`
	PermitWithoutStream bool   `json:"permit_without_stream" yaml:"permit_without_stream"`
}

// Bind binds the actual value from the Servers struct field.
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

	if s.TLS != nil { // This handling is compliant as per previous similar cases
		s.TLS.Bind()
	} else {
		s.TLS = &TLS{
			Enabled: false,
		}
	}
	return s
}

// GetGRPCStreamConcurrency returns the gRPC stream concurrency.
func (s *Servers) GetGRPCStreamConcurrency() (c int) {
	for _, s := range s.Servers {
		if s.GRPC != nil {
			return s.GRPC.BidirectionalStreamConcurrency
		}
	}
	return 0
}

// Bind binds the actual value from the HTTP struct field.
func (h *HTTP) Bind() *HTTP {
	h.HandlerTimeout = GetActualValue(h.HandlerTimeout)
	h.ShutdownDuration = GetActualValue(h.ShutdownDuration)
	h.ReadHeaderTimeout = GetActualValue(h.ReadHeaderTimeout)
	h.ReadTimeout = GetActualValue(h.ReadTimeout)
	h.WriteTimeout = GetActualValue(h.WriteTimeout)
	h.IdleTimeout = GetActualValue(h.IdleTimeout)

	if h.HTTP2 == nil {
		h.HTTP2 = new(HTTP2)
	}
	if h.HTTP2 != nil {
		h.HTTP2.Bind()
	}
	return h
}

// Bind binds the actual value from the GRPC struct field.
func (g *GRPC) Bind() *GRPC {
	g.ConnectionTimeout = GetActualValue(g.ConnectionTimeout)
	g.Interceptors = GetActualValues(g.Interceptors)
	if g.Keepalive == nil {
		g.Keepalive = new(GRPCKeepalive)
	}
	if g.Keepalive != nil {
		g.Keepalive.Bind()
	}
	return g
}

// Bind binds the actual value from the GRPCKeepalive struct field.
func (k *GRPCKeepalive) Bind() *GRPCKeepalive {
	k.MaxConnIdle = GetActualValue(k.MaxConnIdle)
	k.MaxConnAge = GetActualValue(k.MaxConnAge)
	k.MaxConnAgeGrace = GetActualValue(k.MaxConnAgeGrace)
	k.Time = GetActualValue(k.Time)
	k.Timeout = GetActualValue(k.Timeout)
	k.MinTime = GetActualValue(k.MinTime)
	return k
}

// Bind binds the actual value from the Server struct field.
func (s *Server) Bind() *Server {
	s.Name = GetActualValue(s.Name)
	s.Network = GetActualValue(s.Network)
	s.SocketPath = GetActualValue(s.SocketPath)
	s.Host = GetActualValue(s.Host)
	s.Mode = GetActualValue(s.Mode)
	s.ProbeWaitTime = GetActualValue(s.ProbeWaitTime)

	if s.HTTP != nil {
		s.HTTP.Bind()
	}

	if s.GRPC != nil {
		s.GRPC.Bind()
	}
	// Applying the consistent pattern:
	if s.SocketOption == nil {
		s.SocketOption = new(SocketOption)
	}
	s.SocketOption.Bind()
	return s
}

// Opts sets the functional options into the []server.Option slice using the Server struct fields' value.
func (s *Server) Opts() []server.Option {
	opts := make([]server.Option, 0, 10)
	nt := net.NetworkTypeFromString(s.Network)
	if nt == 0 || nt == net.Unknown || strings.EqualFold(nt.String(), net.Unknown.String()) {
		nt = net.TCP
	}
	s.Network = nt.String()
	opts = append(opts,
		server.WithNetwork(s.Network),
		server.WithSocketPath(s.SocketPath),
		server.WithName(s.Name),
		server.WithHost(s.Host),
		server.WithPort(s.Port),
		server.WithProbeWaitTime(s.ProbeWaitTime),
	)

	if s.SocketOption != nil {
		opts = append(opts, server.WithSocketFlag(s.SocketOption.ToSocketFlag()))
	}

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
			if s.HTTP.HTTP2 != nil && s.HTTP.HTTP2.Enabled {
				opts = append(opts,
					server.WithHTTP2Enabled(s.HTTP.HTTP2.Enabled),
					server.WithHandlerLimit(s.HTTP.HTTP2.HandlerLimit),
					server.WithPermitProhibitedCipherSuites(s.HTTP.HTTP2.PermitProhibitedCipherSuites),
					server.WithMaxConcurrentStreams(s.HTTP.HTTP2.MaxConcurrentStreams),
					server.WithMaxUploadBufferPerConnection(s.HTTP.HTTP2.MaxUploadBufferPerConnection),
					server.WithMaxUploadBufferPerStream(s.HTTP.HTTP2.MaxUploadBufferPerStream),
					server.WithMaxDecoderHeaderTableSize(s.HTTP.HTTP2.MaxDecoderHeaderTableSize),
					server.WithMaxEncoderHeaderTableSize(s.HTTP.HTTP2.MaxEncoderHeaderTableSize),
					server.WithMaxReadFrameSize(s.HTTP.HTTP2.MaxReadFrameSize),
				)
			}
		}
	case server.GRPC:
		opts = append(opts,
			server.WithServerMode(mode),
		)
		if s.GRPC != nil {
			opts = append(opts,
				server.WithServerMode(mode),
				server.WithGRPCConnectionTimeout(s.GRPC.ConnectionTimeout),
				server.WithGRPCHeaderTableSize(s.GRPC.HeaderTableSize),
				server.WithGRPCInitialConnWindowSize(s.GRPC.InitialConnWindowSize),
				server.WithGRPCInitialWindowSize(s.GRPC.InitialWindowSize),
				server.WithGRPCInterceptors(s.GRPC.Interceptors...),
				server.WithGRPCMaxConcurrentStreams(s.GRPC.MaxConcurrentStreams),
				server.WithGRPCMaxHeaderListSize(s.GRPC.MaxHeaderListSize),
				server.WithGRPCMaxReceiveMessageSize(s.GRPC.MaxReceiveMessageSize),
				server.WithGRPCMaxSendMessageSize(s.GRPC.MaxSendMessageSize),
				server.WithGRPCNumStreamWorkers(s.GRPC.NumStreamWorkers),
				server.WithGRPCReadBufferSize(s.GRPC.ReadBufferSize),
				server.WithGRPCSharedWriteBuffer(s.GRPC.SharedWriteBuffer),
				server.WithGRPCWaitForHandlers(s.GRPC.WaitForHandlers),
				server.WithGRPCWriteBufferSize(s.GRPC.WriteBufferSize),
			)

			if s.GRPC.EnableReflection {
				opts = append(opts,
					server.WithGRPCRegisterFunc(func(srv *grpc.Server) {
						reflection.Register(srv)
					}))
			}
			if s.GRPC.EnableAdmin || s.GRPC.EnableChannelz {
				opts = append(opts,
					server.WithGRPCRegisterFunc(func(srv *grpc.Server) {
						admin.Register(srv)
					}))
			}
			if s.GRPC.Keepalive != nil {
				opts = append(opts,
					server.WithGRPCKeepaliveMaxConnIdle(s.GRPC.Keepalive.MaxConnIdle),
					server.WithGRPCKeepaliveMaxConnAge(s.GRPC.Keepalive.MaxConnAge),
					server.WithGRPCKeepaliveMaxConnAgeGrace(s.GRPC.Keepalive.MaxConnAgeGrace),
					server.WithGRPCKeepaliveTime(s.GRPC.Keepalive.Time),
					server.WithGRPCKeepaliveTimeout(s.GRPC.Keepalive.Timeout),
					server.WithGRPCKeepaliveMinTime(s.GRPC.Keepalive.MinTime),
					server.WithGRPCKeepalivePermitWithoutStream(s.GRPC.Keepalive.PermitWithoutStream),
				)
			}
		}
	}

	return opts
}
