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
	// GRPC represents the gRPC configuration.
	GRPC *GRPC `json:"grpc,omitempty" yaml:"grpc"`
	// SocketOption represents the socket option configuration.
	SocketOption *SocketOption `json:"socket_option,omitempty" yaml:"socket_option"`
	// HTTP represents the HTTP configuration.
	HTTP *HTTP `json:"http,omitempty" yaml:"http"`
	// Name represents the server name.
	Name string `json:"name,omitempty" yaml:"name"`
	// Network represents the network type (tcp, unix, etc.).
	Network string `json:"network,omitempty" yaml:"network"`
	// Host represents the server host.
	Host string `json:"host,omitempty" yaml:"host"`
	// SocketPath represents the socket path.
	SocketPath string `json:"socket_path,omitempty" yaml:"socket_path"`
	// Mode represents the server mode (gRPC, REST, GraphQL).
	Mode string `json:"mode,omitempty" yaml:"mode"` // gRPC, REST, GraphQL
	// ProbeWaitTime represents the probe wait time.
	ProbeWaitTime string `json:"probe_wait_time,omitempty" yaml:"probe_wait_time"`
	// Restart enables server restart on failure.
	Restart bool `json:"restart,omitempty" yaml:"restart"`
	// Port represents the server port.
	Port uint16 `json:"port,omitempty" yaml:"port"`
}

// HTTP represents the configuration for HTTP.
type HTTP struct {
	// HTTP2 represents the configuration for HTTP2.
	HTTP2 *HTTP2 `json:"http2" yaml:"http2"`
	// ShutdownDuration represents the duration for the http server to shutdown.
	ShutdownDuration string `json:"shutdown_duration" yaml:"shutdown_duration"`
	// HandlerTimeout represents the timeout duration for http handlers.
	HandlerTimeout string `json:"handler_timeout" yaml:"handler_timeout"`
	// IdleTimeout represents the maximum amount of time to wait for the next request when keep-alives are enabled.
	IdleTimeout string `json:"idle_timeout" yaml:"idle_timeout"`
	// ReadHeaderTimeout represents the amount of time allowed to read request headers.
	ReadHeaderTimeout string `json:"read_header_timeout" yaml:"read_header_timeout"`
	// ReadTimeout represents the maximum duration for reading the entire request, including the body.
	ReadTimeout string `json:"read_timeout" yaml:"read_timeout"`
	// WriteTimeout represents the maximum duration before timing out writes of the response.
	WriteTimeout string `json:"write_timeout" yaml:"write_timeout"`
}

// HTTP2 represents the configuration for HTTP2.
type HTTP2 struct {
	// HandlerLimit represents the limit of handlers.
	HandlerLimit int `json:"handler_limit,omitempty" yaml:"handler_limit"`
	// Enabled enables HTTP2.
	Enabled bool `json:"enabled,omitempty" yaml:"enabled"`
	// PermitProhibitedCipherSuites enables prohibited cipher suites.
	PermitProhibitedCipherSuites bool `json:"permit_prohibited_cipher_suites,omitempty" yaml:"permit_prohibited_cipher_suites"`
	// MaxUploadBufferPerConnection represents the maximum upload buffer per connection.
	MaxUploadBufferPerConnection int32 `json:"max_upload_buffer_per_connection,omitempty" yaml:"max_upload_buffer_per_connection"`
	// MaxUploadBufferPerStream represents the maximum upload buffer per stream.
	MaxUploadBufferPerStream int32 `json:"max_upload_buffer_per_stream,omitempty" yaml:"max_upload_buffer_per_stream"`
	// MaxConcurrentStreams represents the maximum concurrent streams.
	MaxConcurrentStreams uint32 `json:"max_concurrent_streams,omitempty" yaml:"max_concurrent_streams"`
	// MaxDecoderHeaderTableSize represents the maximum decoder header table size.
	MaxDecoderHeaderTableSize uint32 `json:"max_decoder_header_table_size,omitempty" yaml:"max_decoder_header_table_size"`
	// MaxEncoderHeaderTableSize represents the maximum encoder header table size.
	MaxEncoderHeaderTableSize uint32 `json:"max_encoder_header_table_size,omitempty" yaml:"max_encoder_header_table_size"`
	// MaxReadFrameSize represents the maximum read frame size.
	MaxReadFrameSize uint32 `json:"max_read_frame_size,omitempty" yaml:"max_read_frame_size"`
}

// Bind binds the actual data from the HTTP2 receiver fields.
func (h *HTTP2) Bind() *HTTP2 {
	// No fields to bind as per rules
	return h
}

// GRPC represents the configuration for gPRC.
type GRPC struct {
	// Keepalive represents the gRPC keepalive configuration.
	Keepalive *GRPCKeepalive `json:"keepalive,omitempty" yaml:"keepalive"`
	// ConnectionTimeout represents the gRPC connection timeout duration.
	ConnectionTimeout string `json:"connection_timeout,omitempty" yaml:"connection_timeout"`
	// Interceptors represents the list of gRPC interceptors.
	Interceptors []string `json:"interceptors,omitempty" yaml:"interceptors"`
	// InitialConnWindowSize represents the initial connection window size.
	InitialConnWindowSize int `json:"initial_conn_window_size,omitempty" yaml:"initial_conn_window_size"`
	// WriteBufferSize represents the write buffer size.
	WriteBufferSize int `json:"write_buffer_size,omitempty" yaml:"write_buffer_size"`
	// ReadBufferSize represents the read buffer size.
	ReadBufferSize int `json:"read_buffer_size,omitempty" yaml:"read_buffer_size"`
	// MaxSendMessageSize represents the maximum send message size.
	MaxSendMessageSize int `json:"max_send_message_size,omitempty" yaml:"max_send_message_size"`
	// MaxReceiveMessageSize represents the maximum receive message size.
	MaxReceiveMessageSize int `json:"max_receive_message_size,omitempty" yaml:"max_receive_message_size"`
	// InitialWindowSize represents the initial window size.
	InitialWindowSize int `json:"initial_window_size,omitempty" yaml:"initial_window_size"`
	// BidirectionalStreamConcurrency represents the bidirectional stream concurrency limit.
	BidirectionalStreamConcurrency int `json:"bidirectional_stream_concurrency,omitempty" yaml:"bidirectional_stream_concurrency"`
	// NumStreamWorkers represents the number of stream workers.
	NumStreamWorkers uint32 `json:"num_stream_workers,omitempty" yaml:"num_stream_workers"`
	// MaxHeaderListSize represents the maximum header list size.
	MaxHeaderListSize uint32 `json:"max_header_list_size,omitempty" yaml:"max_header_list_size"`
	// MaxConcurrentStreams represents the maximum concurrent streams.
	MaxConcurrentStreams uint32 `json:"max_concurrent_streams,omitempty" yaml:"max_concurrent_streams"`
	// HeaderTableSize represents the header table size.
	HeaderTableSize uint32 `json:"header_table_size,omitempty" yaml:"header_table_size"`
	// EnableAdmin enables the admin service.
	EnableAdmin bool `json:"enable_admin,omitempty" yaml:"enable_admin"`
	// WaitForHandlers waits for handlers to finish.
	WaitForHandlers bool `json:"wait_for_handlers,omitempty" yaml:"wait_for_handlers"`
	// SharedWriteBuffer enables the shared write buffer.
	SharedWriteBuffer bool `json:"shared_write_buffer,omitempty" yaml:"shared_write_buffer"`
	// EnableReflection enables the reflection service.
	EnableReflection bool `json:"enable_reflection,omitempty" yaml:"enable_reflection"`
	// EnableChannelz enables the channelz service.
	EnableChannelz bool `json:"enable_channelz,omitempty" yaml:"enable_channelz"`
}

// GRPCKeepalive represents the configuration for gRPC keep-alive.
type GRPCKeepalive struct {
	// MaxConnIdle represents the maximum amount of time a connection may be idle.
	MaxConnIdle string `json:"max_conn_idle" yaml:"max_conn_idle"`
	// MaxConnAge represents the maximum amount of time a connection may exist.
	MaxConnAge string `json:"max_conn_age" yaml:"max_conn_age"`
	// MaxConnAgeGrace represents the additive period after MaxConnAge after which the connection will be forcibly closed.
	MaxConnAgeGrace string `json:"max_conn_age_grace" yaml:"max_conn_age_grace"`
	// Time represents the duration after which if the client doesn't see any activity it pings the server to see if the transport is still alive.
	Time string `json:"time" yaml:"time"`
	// Timeout represents the duration that the client waits for a ping response.
	Timeout string `json:"timeout" yaml:"timeout"`
	// MinTime represents the minimum amount of time a client should wait before sending a keepalive ping.
	MinTime string `json:"min_time" yaml:"min_time"`
	// PermitWithoutStream if true, client can send keepalive pings even with no active RPCs.
	PermitWithoutStream bool `json:"permit_without_stream" yaml:"permit_without_stream"`
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
					server.WithGRPCRegisterar(func(srv *grpc.Server) {
						reflection.Register(srv)
					}))
			}
			if s.GRPC.EnableAdmin || s.GRPC.EnableChannelz {
				opts = append(opts,
					server.WithGRPCRegisterar(func(srv *grpc.Server) {
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
