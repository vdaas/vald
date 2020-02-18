//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/tls"
)

type GRPCClient struct {
	Addrs               []string    `json:"addrs" yaml:"addrs"`
	HealthCheckDuration string      `json:"health_check_duration" yaml:"health_check_duration"`
	ConnectionPool      int         `json:"connection_pool" yaml:"connection_pool"`
	Backoff             *Backoff    `json:"backoff" yaml:"backoff"`
	CallOption          *CallOption `json:"call_option" yaml:"call_option"`
	DialOption          *DialOption `json:"dial_option" yaml:"dial_option"`
	TLS                 *TLS        `json:"tls" yaml:"tls"`
}

type CallOption struct {
	WaitForReady          bool `json:"wait_for_ready" yaml:"wait_for_ready"`
	MaxRetryRPCBufferSize int  `json:"max_retry_rpc_buffer_size" yaml:"max_retry_rpc_buffer_size"`
	MaxRecvMsgSize        int  `json:"max_recv_msg_size" yaml:"max_recv_msg_size"`
	MaxSendMsgSize        int  `json:"max_send_msg_size" yaml:"max_send_msg_size"`
}

type DialOption struct {
	WriteBufferSize             int                  `json:"write_buffer_size" yaml:"write_buffer_size"`
	ReadBufferSize              int                  `json:"read_buffer_size" yaml:"read_buffer_size"`
	InitialWindowSize           int                  `json:"initial_window_size" yaml:"initial_window_size"`
	InitialConnectionWindowSize int                  `json:"initial_connection_window_size" yaml:"initial_connection_window_size"`
	MaxMsgSize                  int                  `json:"max_msg_size" yaml:"max_msg_size"`
	MaxBackoffDelay             string               `json:"max_backoff_delay" yaml:"max_backoff_delay"`
	EnableBackoff               bool                 `json:"enable_backoff" yaml:"enable_backoff"`
	Insecure                    bool                 `json:"insecure" yaml:"insecure"`
	Timeout                     string               `json:"timeout" yaml:"timeout"`
	TCP                         *TCP                 `json:"tcp" yaml:"tcp"`
	KeepAlive                   *GRPCClientKeepalive `json:"keep_alive" yaml:"keep_alive"`
}

type GRPCClientKeepalive struct {
	Time                string `json:"time" yaml:"time"`
	Timeout             string `json:"timeout" yaml:"timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream" yaml:"permit_without_stream"`
}

func newGRPCClientConfig() *GRPCClient {
	return &GRPCClient{
		DialOption: &DialOption{
			Insecure: true,
		},
	}
}

func (g *GRPCClient) Bind() *GRPCClient {
	g.HealthCheckDuration = GetActualValue(g.HealthCheckDuration)

	g.Addrs = GetActualValues(g.Addrs)

	if g.Backoff != nil {
		g.Backoff.Bind()
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

func (g *GRPCClientKeepalive) Bind() *GRPCClientKeepalive {
	g.Time = GetActualValue(g.Time)
	g.Timeout = GetActualValue(g.Timeout)
	return g
}

func (c *CallOption) Bind() *CallOption {
	return c
}

func (d *DialOption) Bind() *DialOption {
	d.MaxBackoffDelay = GetActualValue(d.MaxBackoffDelay)
	d.Timeout = GetActualValue(d.Timeout)
	return d
}

func (g *GRPCClient) Opts() []grpc.Option {
	opts := make([]grpc.Option, 0, 18)
	opts = append(opts,
		grpc.WithHealthCheckDuration(g.HealthCheckDuration),
		grpc.WithConnectionPool(g.ConnectionPool),
	)
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
	if g.CallOption != nil {
		opts = append(opts,
			grpc.WithWaitForReady(g.CallOption.WaitForReady),
			grpc.WithMaxRetryRPCBufferSize(g.CallOption.MaxRetryRPCBufferSize),
			grpc.WithMaxRecvMsgSize(g.CallOption.MaxRecvMsgSize),
			grpc.WithMaxSendMsgSize(g.CallOption.MaxSendMsgSize),
		)
	}

	if g.DialOption != nil {
		opts = append(opts,
			grpc.WithWriteBufferSize(g.DialOption.WriteBufferSize),
			grpc.WithReadBufferSize(g.DialOption.WriteBufferSize),
			grpc.WithInitialWindowSize(g.DialOption.InitialWindowSize),
			grpc.WithInitialConnectionWindowSize(g.DialOption.InitialWindowSize),
			grpc.WithMaxMsgSize(g.DialOption.MaxMsgSize),
			grpc.WithInsecure(g.DialOption.Insecure),
			grpc.WithMaxBackoffDelay(g.DialOption.MaxBackoffDelay),
			grpc.WithDialTimeout(g.DialOption.Timeout),
		)

		if g.DialOption.TCP != nil &&
			len(g.DialOption.TCP.Dialer.Timeout) != 0 {
			if g.DialOption.TCP.TLS != nil && g.DialOption.TCP.TLS.Enabled {
				opts = append(opts,
					grpc.WithInsecure(false),
				)
			}
			der, err := tcp.NewDialer(g.DialOption.TCP.Opts()...)
			if err == nil {
				opts = append(opts,
					grpc.WithDialer(der),
				)
			}
		}

		if g.DialOption.KeepAlive != nil {
			opts = append(opts,
				grpc.WithKeepaliveParams(
					g.DialOption.KeepAlive.Time,
					g.DialOption.KeepAlive.Timeout,
					g.DialOption.KeepAlive.PermitWithoutStream,
				),
			)
		}
	}

	if g.TLS != nil && g.TLS.Enabled {
		cfg, err := tls.NewClientConfig(g.TLS.Opts()...)
		if err == nil && cfg != nil {
			opts = append(opts,
				grpc.WithTLSConfig(cfg),
			)
		}
	}

	return opts
}
