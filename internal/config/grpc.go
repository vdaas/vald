//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/tls"
)

// GRPCClient represents the configurations for gRPC client.
type GRPCClient struct {
	Addrs               []string        `json:"addrs"                 yaml:"addrs"`
	HealthCheckDuration string          `json:"health_check_duration" yaml:"health_check_duration"`
	ConnectionPool      *ConnectionPool `json:"connection_pool"       yaml:"connection_pool"`
	Backoff             *Backoff        `json:"backoff"               yaml:"backoff"`
	CircuitBreaker      *CircuitBreaker `json:"circuit_breaker"       yaml:"circuit_breaker"`
	CallOption          *CallOption     `json:"call_option"           yaml:"call_option"`
	DialOption          *DialOption     `json:"dial_option"           yaml:"dial_option"`
	TLS                 *TLS            `json:"tls"                   yaml:"tls"`
}

// CallOption represents the configurations for call option.
type CallOption struct {
	WaitForReady          bool   `json:"wait_for_ready"            yaml:"wait_for_ready"`
	MaxRetryRPCBufferSize int    `json:"max_retry_rpc_buffer_size" yaml:"max_retry_rpc_buffer_size"`
	MaxRecvMsgSize        int    `json:"max_recv_msg_size"         yaml:"max_recv_msg_size"`
	MaxSendMsgSize        int    `json:"max_send_msg_size"         yaml:"max_send_msg_size"`
	ContentSubtype        string `json:"content_subtype"           yaml:"content_subtype"`
}

// DialOption represents the configurations for dial option.
type DialOption struct {
	DisableRetry                bool                 `json:"disable_retry,omitempty"                  yaml:"disable_retry"`
	EnableBackoff               bool                 `json:"enable_backoff,omitempty"                 yaml:"enable_backoff"`
	Insecure                    bool                 `json:"insecure,omitempty"                       yaml:"insecure"`
	SharedWriteBuffer           bool                 `json:"shared_write_buffer,omitempty"            yaml:"shared_write_buffer"`
	InitialConnectionWindowSize int32                `json:"initial_connection_window_size,omitempty" yaml:"initial_connection_window_size"`
	InitialWindowSize           int32                `json:"initial_window_size,omitempty"            yaml:"initial_window_size"`
	MaxHeaderListSize           uint32               `json:"max_header_list_size,omitempty"           yaml:"max_header_list_size"`
	MaxCallAttempts             int                  `json:"max_call_attempts,omitempty"              yaml:"max_call_attempts"`
	MaxMsgSize                  int                  `json:"max_msg_size,omitempty"                   yaml:"max_msg_size"`
	ReadBufferSize              int                  `json:"read_buffer_size,omitempty"               yaml:"read_buffer_size"`
	WriteBufferSize             int                  `json:"write_buffer_size,omitempty"              yaml:"write_buffer_size"`
	BackoffJitter               float64              `json:"backoff_jitter,omitempty"                 yaml:"backoff_jitter"`
	BackoffMultiplier           float64              `json:"backoff_multiplier,omitempty"             yaml:"backoff_multiplier"`
	Authority                   string               `json:"authority,omitempty"                      yaml:"authority"`
	BackoffBaseDelay            string               `json:"backoff_base_delay,omitempty"             yaml:"backoff_base_delay"`
	BackoffMaxDelay             string               `json:"backoff_max_delay,omitempty"              yaml:"backoff_max_delay"`
	IdleTimeout                 string               `json:"idle_timeout,omitempty"                   yaml:"idle_timeout"`
	MinimumConnectionTimeout    string               `json:"min_connection_timeout,omitempty"         yaml:"min_connection_timeout"`
	Timeout                     string               `json:"timeout,omitempty"                        yaml:"timeout"`
	UserAgent                   string               `json:"user_agent,omitempty"                     yaml:"user_agent"`
	Interceptors                []string             `json:"interceptors,omitempty"                   yaml:"interceptors"`
	Net                         *Net                 `json:"net,omitempty"                            yaml:"net"`
	Keepalive                   *GRPCClientKeepalive `json:"keepalive,omitempty"                      yaml:"keepalive"`
}

// ConnectionPool represents the configurations for connection pool.
type ConnectionPool struct {
	ResolveDNS           bool   `json:"enable_dns_resolver"     yaml:"enable_dns_resolver"`
	EnableRebalance      bool   `json:"enable_rebalance"        yaml:"enable_rebalance"`
	RebalanceDuration    string `json:"rebalance_duration"      yaml:"rebalance_duration"`
	Size                 int    `json:"size"                    yaml:"size"`
	OldConnCloseDuration string `json:"old_conn_close_duration" yaml:"old_conn_close_duration"`
}

// GRPCClientKeepalive represents the configurations for gRPC keep-alive.
type GRPCClientKeepalive struct {
	Time                string `json:"time"                  yaml:"time"`
	Timeout             string `json:"timeout"               yaml:"timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream" yaml:"permit_without_stream"`
}

// newGRPCClientConfig returns the GRPCClient with DailOption with insecure is true.
func newGRPCClientConfig() *GRPCClient {
	return &GRPCClient{
		DialOption: &DialOption{
			Insecure: true,
		},
	}
}

// Bind binds the actual data from the GRPCClient receiver fields.
func (g *GRPCClient) Bind() *GRPCClient {
	g.Addrs = GetActualValues(g.Addrs)
	g.HealthCheckDuration = GetActualValue(g.HealthCheckDuration)

	if g.ConnectionPool != nil {
		g.ConnectionPool.RebalanceDuration = GetActualValue(g.ConnectionPool.RebalanceDuration)
		g.ConnectionPool.OldConnCloseDuration = GetActualValue(g.ConnectionPool.OldConnCloseDuration)
	} else {
		g.ConnectionPool = new(ConnectionPool)
	}

	if g.Backoff != nil {
		g.Backoff.Bind()
	}

	if g.CircuitBreaker != nil {
		g.CircuitBreaker.Bind()
	}

	if g.CallOption != nil {
		g.CallOption.Bind()
	}

	if g.DialOption != nil {
		g.DialOption.Bind()
	} else {
		g.DialOption = new(DialOption)
	}

	if g.TLS != nil &&
		g.TLS.Enabled &&
		g.TLS.Cert != "" &&
		g.TLS.Key != "" {
		g.TLS.Bind()
	} else {
		g.TLS = &TLS{
			Enabled: false,
		}
		g.DialOption.Insecure = true
	}

	return g
}

// Bind binds the actual data from the GRPCClientKeepalive receiver fields.
func (g *GRPCClientKeepalive) Bind() *GRPCClientKeepalive {
	g.Time = GetActualValue(g.Time)
	g.Timeout = GetActualValue(g.Timeout)
	return g
}

// Bind binds the actual data from the CallOption receiver fields.
func (c *CallOption) Bind() *CallOption {
	c.ContentSubtype = GetActualValue(c.ContentSubtype)
	return c
}

// Bind binds the actual data from the DialOption receiver fields.
func (d *DialOption) Bind() *DialOption {
	d.Authority = GetActualValue(d.Authority)
	d.BackoffBaseDelay = GetActualValue(d.BackoffBaseDelay)
	d.BackoffMaxDelay = GetActualValue(d.BackoffMaxDelay)
	d.IdleTimeout = GetActualValue(d.IdleTimeout)
	d.Interceptors = GetActualValues(d.Interceptors)
	d.MinimumConnectionTimeout = GetActualValue(d.MinimumConnectionTimeout)
	d.Timeout = GetActualValue(d.Timeout)
	d.UserAgent = GetActualValue(d.UserAgent)
	return d
}

// Opts creates the slice with the functional options for the gRPC options.
func (g *GRPCClient) Opts() ([]grpc.Option, error) {
	if g == nil {
		return nil, nil
	}
	opts := make([]grpc.Option, 0, 18)

	if len(g.HealthCheckDuration) != 0 {
		opts = append(opts, grpc.WithHealthCheckDuration(g.HealthCheckDuration))
	}

	if g.ConnectionPool != nil {
		opts = append(opts,
			grpc.WithConnectionPoolSize(g.ConnectionPool.Size),
			grpc.WithOldConnCloseDelay(g.ConnectionPool.OldConnCloseDuration),
			grpc.WithResolveDNS(g.ConnectionPool.ResolveDNS),
		)
		if g.ConnectionPool.EnableRebalance {
			opts = append(opts,
				grpc.WithEnableConnectionPoolRebalance(g.ConnectionPool.EnableRebalance),
				grpc.WithConnectionPoolRebalanceDuration(g.ConnectionPool.RebalanceDuration),
			)
		}
	}

	if g.Addrs != nil && len(g.Addrs) != 0 {
		opts = append(opts,
			grpc.WithAddrs(g.Addrs...),
		)
	}

	if g.Backoff != nil &&
		len(g.Backoff.InitialDuration) != 0 &&
		g.Backoff.RetryCount > 2 {
		opts = append(opts,
			grpc.WithBackoff(
				backoff.New(g.Backoff.Opts()...),
			),
		)
	}

	if g.CircuitBreaker != nil {
		cb, err := circuitbreaker.NewCircuitBreaker(
			circuitbreaker.WithBreakerOpts(
				circuitbreaker.WithClosedErrorRate(g.CircuitBreaker.ClosedErrorRate),
				circuitbreaker.WithHalfOpenErrorRate(g.CircuitBreaker.HalfOpenErrorRate),
				circuitbreaker.WithMinSamples(g.CircuitBreaker.MinSamples),
				circuitbreaker.WithOpenTimeout(g.CircuitBreaker.OpenTimeout),
				circuitbreaker.WithClosedRefreshTimeout(g.CircuitBreaker.ClosedRefreshTimeout),
			),
		)
		if err != nil {
			return nil, err
		}
		opts = append(opts,
			grpc.WithCircuitBreaker(cb),
		)
	}

	if g.CallOption != nil {
		opts = append(opts,
			grpc.WithCallContentSubtype(g.CallOption.ContentSubtype),
			grpc.WithMaxRecvMsgSize(g.CallOption.MaxRecvMsgSize),
			grpc.WithMaxRetryRPCBufferSize(g.CallOption.MaxRetryRPCBufferSize),
			grpc.WithMaxSendMsgSize(g.CallOption.MaxSendMsgSize),
			grpc.WithWaitForReady(g.CallOption.WaitForReady),
		)
	}

	if g.DialOption != nil {
		opts = append(opts,
			grpc.WithAuthority(g.DialOption.Authority),
			grpc.WithBackoffMaxDelay(g.DialOption.BackoffMaxDelay),
			grpc.WithClientInterceptors(g.DialOption.Interceptors...),
			grpc.WithDisableRetry(g.DialOption.DisableRetry),
			grpc.WithIdleTimeout(g.DialOption.IdleTimeout),
			grpc.WithInitialConnectionWindowSize(g.DialOption.InitialConnectionWindowSize),
			grpc.WithInitialWindowSize(g.DialOption.InitialWindowSize),
			grpc.WithInsecure(g.DialOption.Insecure),
			grpc.WithMaxCallAttempts(g.DialOption.MaxCallAttempts),
			grpc.WithMaxHeaderListSize(g.DialOption.MaxHeaderListSize),
			grpc.WithMaxMsgSize(g.DialOption.MaxMsgSize),
			grpc.WithReadBufferSize(g.DialOption.ReadBufferSize),
			grpc.WithSharedWriteBuffer(g.DialOption.SharedWriteBuffer),
			grpc.WithUserAgent(g.DialOption.UserAgent),
			grpc.WithWriteBufferSize(g.DialOption.WriteBufferSize),
		)

		if g.DialOption.Net != nil && g.DialOption.Net.Dialer != nil &&
			len(g.DialOption.Net.Dialer.Timeout) != 0 {
			if g.DialOption.Net.TLS != nil && g.DialOption.Net.TLS.Enabled {
				opts = append(opts,
					grpc.WithInsecure(false),
				)
			}
			netOpts, err := g.DialOption.Net.Opts()
			if err != nil {
				return nil, err
			}
			der, err := net.NewDialer(netOpts...)
			if err != nil {
				return nil, err
			}
			network := g.DialOption.Net.Network
			if network == "" {
				network = net.TCP.String()
			}
			opts = append(opts,
				grpc.WithDialer(network, der),
			)
		}

		if g.DialOption.Keepalive != nil {
			opts = append(opts,
				grpc.WithKeepaliveParams(
					g.DialOption.Keepalive.Time,
					g.DialOption.Keepalive.Timeout,
					g.DialOption.Keepalive.PermitWithoutStream,
				),
			)
		}
	}

	if g.TLS != nil && g.TLS.Enabled {
		cfg, err := tls.NewClientConfig(g.TLS.Opts()...)
		if err != nil {
			return nil, err
		}
		opts = append(opts,
			grpc.WithTLSConfig(cfg),
		)
	}

	return opts, nil
}
