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
	"slices"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/tls"
)

// GRPCClient represents the configurations for gRPC client.
type GRPCClient struct {
	// ConnectionPool represents the connection pool configuration.
	ConnectionPool *ConnectionPool `json:"connection_pool" yaml:"connection_pool"`
	// Backoff represents the backoff configuration.
	Backoff *Backoff `json:"backoff" yaml:"backoff"`
	// CircuitBreaker represents the circuit breaker configuration.
	CircuitBreaker *CircuitBreaker `json:"circuit_breaker" yaml:"circuit_breaker"`
	// CallOption represents the call option configuration.
	CallOption *CallOption `json:"call_option" yaml:"call_option"`
	// DialOption represents the dial option configuration.
	DialOption *DialOption `json:"dial_option" yaml:"dial_option"`
	// TLS represents the TLS configuration.
	TLS *TLS `json:"tls" yaml:"tls"`
	// HealthCheckDuration represents the health check duration.
	HealthCheckDuration string `json:"health_check_duration" yaml:"health_check_duration"`
	// Addrs represents the list of addresses.
	Addrs []string `json:"addrs" yaml:"addrs"`
}

// CallOption represents the configurations for call option.
type CallOption struct {
	// ContentSubtype represents the content subtype.
	ContentSubtype string `json:"content_subtype" yaml:"content_subtype"`
	// MaxRetryRPCBufferSize represents the maximum retry RPC buffer size.
	MaxRetryRPCBufferSize int `json:"max_retry_rpc_buffer_size" yaml:"max_retry_rpc_buffer_size"`
	// MaxRecvMsgSize represents the maximum receive message size.
	MaxRecvMsgSize int `json:"max_recv_msg_size" yaml:"max_recv_msg_size"`
	// MaxSendMsgSize represents the maximum send message size.
	MaxSendMsgSize int `json:"max_send_msg_size" yaml:"max_send_msg_size"`
	// WaitForReady enables wait for ready.
	WaitForReady bool `json:"wait_for_ready" yaml:"wait_for_ready"`
}

// DialOption represents the configurations for dial option.
type DialOption struct {
	// Keepalive represents the keepalive configuration.
	Keepalive *GRPCClientKeepalive `json:"keepalive,omitempty" yaml:"keepalive"`
	// Net represents the network configuration.
	Net *Net `json:"net,omitempty" yaml:"net"`
	// Authority represents the authority.
	Authority string `json:"authority,omitempty" yaml:"authority"`
	// UserAgent represents the user agent.
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent"`
	// Timeout represents the timeout duration.
	Timeout string `json:"timeout,omitempty" yaml:"timeout"`
	// MinimumConnectionTimeout represents the minimum connection timeout.
	MinimumConnectionTimeout string `json:"min_connection_timeout,omitempty" yaml:"min_connection_timeout"`
	// IdleTimeout represents the idle timeout.
	IdleTimeout string `json:"idle_timeout,omitempty" yaml:"idle_timeout"`
	// BackoffMaxDelay represents the backoff maximum delay.
	BackoffMaxDelay string `json:"backoff_max_delay,omitempty" yaml:"backoff_max_delay"`
	// BackoffBaseDelay represents the backoff base delay.
	BackoffBaseDelay string `json:"backoff_base_delay,omitempty" yaml:"backoff_base_delay"`
	// Interceptors represents the interceptors.
	Interceptors []string `json:"interceptors,omitempty" yaml:"interceptors"`
	// WriteBufferSize represents the write buffer size.
	WriteBufferSize int `json:"write_buffer_size,omitempty" yaml:"write_buffer_size"`
	// BackoffJitter represents the backoff jitter.
	BackoffJitter float64 `json:"backoff_jitter,omitempty" yaml:"backoff_jitter"`
	// BackoffMultiplier represents the backoff multiplier.
	BackoffMultiplier float64 `json:"backoff_multiplier,omitempty" yaml:"backoff_multiplier"`
	// ReadBufferSize represents the read buffer size.
	ReadBufferSize int `json:"read_buffer_size,omitempty" yaml:"read_buffer_size"`
	// MaxMsgSize represents the maximum message size.
	MaxMsgSize int `json:"max_msg_size,omitempty" yaml:"max_msg_size"`
	// MaxCallAttempts represents the maximum call attempts.
	MaxCallAttempts int `json:"max_call_attempts,omitempty" yaml:"max_call_attempts"`
	// MaxHeaderListSize represents the maximum header list size.
	MaxHeaderListSize uint32 `json:"max_header_list_size,omitempty" yaml:"max_header_list_size"`
	// InitialWindowSize represents the initial window size.
	InitialWindowSize int32 `json:"initial_window_size,omitempty" yaml:"initial_window_size"`
	// InitialConnectionWindowSize represents the initial connection window size.
	InitialConnectionWindowSize int32 `json:"initial_connection_window_size,omitempty" yaml:"initial_connection_window_size"`
	// DisableRetry disables retry.
	DisableRetry bool `json:"disable_retry,omitempty" yaml:"disable_retry"`
	// SharedWriteBuffer enables shared write buffer.
	SharedWriteBuffer bool `json:"shared_write_buffer,omitempty" yaml:"shared_write_buffer"`
	// Insecure enables insecure mode.
	Insecure bool `json:"insecure,omitempty" yaml:"insecure"`
	// EnableBackoff enables backoff.
	EnableBackoff bool `json:"enable_backoff,omitempty" yaml:"enable_backoff"`
}

// ConnectionPool represents the configurations for connection pool.
type ConnectionPool struct {
	// RebalanceDuration represents the rebalance duration.
	RebalanceDuration string `json:"rebalance_duration" yaml:"rebalance_duration"`
	// OldConnCloseDuration represents the old connection close duration.
	OldConnCloseDuration string `json:"old_conn_close_duration" yaml:"old_conn_close_duration"`
	// Size represents the pool size.
	Size int `json:"size" yaml:"size"`
	// ResolveDNS enables DNS resolution.
	ResolveDNS bool `json:"enable_dns_resolver" yaml:"enable_dns_resolver"`
	// EnableRebalance enables rebalance.
	EnableRebalance bool `json:"enable_rebalance" yaml:"enable_rebalance"`
}

// Bind binds the actual data from the ConnectionPool receiver fields.
func (cp *ConnectionPool) Bind() *ConnectionPool {
	cp.RebalanceDuration = GetActualValue(cp.RebalanceDuration)
	cp.OldConnCloseDuration = GetActualValue(cp.OldConnCloseDuration)
	return cp
}

// GRPCClientKeepalive represents the configurations for gRPC keep-alive.
type GRPCClientKeepalive struct {
	// Time represents the time duration.
	Time string `json:"time" yaml:"time"`
	// Timeout represents the timeout duration.
	Timeout string `json:"timeout" yaml:"timeout"`
	// PermitWithoutStream enables permit without stream.
	PermitWithoutStream bool `json:"permit_without_stream" yaml:"permit_without_stream"`
}

// newGRPCClientConfig returns the GRPCClient with DailOption with insecure is true.
func newGRPCClientConfig() *GRPCClient {
	return (&GRPCClient{
		DialOption: &DialOption{
			Insecure: true,
		},
	}).Bind()
}

// Bind binds the actual data from the GRPCClient receiver fields.
func (g *GRPCClient) Bind() *GRPCClient {
	g.Addrs = GetActualValues(g.Addrs)
	slices.Sort(g.Addrs)
	g.Addrs = slices.Compact(g.Addrs)

	g.HealthCheckDuration = GetActualValue(g.HealthCheckDuration)

	if g.ConnectionPool != nil {
		g.ConnectionPool.Bind()
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

	if g.DialOption != nil { // This part is already compliant due to else clause
		g.DialOption.Bind()
	}

	if g.TLS != nil && g.TLS.Enabled {
		g.TLS.Bind()
	} else {
		g.TLS = (&TLS{
			Enabled:            false,
			InsecureSkipVerify: true,
		}).Bind()
		if g.DialOption != nil {
			g.DialOption.Insecure = true
		}
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

	if d.Net != nil {
		d.Net.Bind()
	}

	if d.Keepalive != nil {
		d.Keepalive.Bind()
	}
	return d
}

// Opts creates the slice with the functional options for the gRPC options.
func (g *GRPCClient) Opts() ([]grpc.Option, error) {
	if g == nil {
		return nil, nil
	}
	opts := make([]grpc.Option, 0, 18)

	if g.HealthCheckDuration != "" {
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

	if len(g.Addrs) != 0 {
		opts = append(opts,
			grpc.WithAddrs(g.Addrs...),
		)
	}

	if g.Backoff != nil &&
		g.Backoff.InitialDuration != "" &&
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
			g.DialOption.Net.Dialer.Timeout != "" {
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
