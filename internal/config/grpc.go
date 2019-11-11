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

type GRPCClient struct {
	HealthCheckDuration string      `json:"health_check_duration" yaml:"health_check_duration"`
	CallOption          *CallOption `json:"call_option" yaml:"call_option"`
	DialOption          *DialOption `json:"dial_option" yaml:"dial_option"`
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
	EnableBackofff              bool                 `json:"enable_backofff" yaml:"enable_backofff"`
	Insecure                    bool                 `json:"insecure" yaml:"insecure"`
	Timeout                     string               `json:"timeout" yaml:"timeout"`
	Dialer                      *TCP                 `json:"dialer" yaml:"dialer"`
	KeepAlive                   *GRPCClientKeepalive `json:"keep_alive" yaml:"keep_alive"`
}

type GRPCClientKeepalive struct {
	MinPingTime         string `json:"min_ping_time" yaml:"min_ping_time"`
	Timeout             string `json:"timeout" yaml:"timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream" yaml:"permit_without_stream"`
}

func (g *GRPCClient) Bind() *GRPCClient {
	g.HealthCheckDuration = GetActualValue(g.HealthCheckDuration)

	if g.CallOption != nil {
		g.CallOption.Bind()
	}

	if g.DialOption != nil {
		g.DialOption.Bind()
	}
	return g
}

func (g *GRPCClientKeepalive) Bind() *GRPCClientKeepalive {
	g.MinPingTime = GetActualValue(g.MinPingTime)
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
