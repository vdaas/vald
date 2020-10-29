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

// Package grpc provides vald gRPC client functions
package grpc

import "github.com/vdaas/vald/internal/config"

// Option is gatewayClient configure.
type Option func(*gatewayClient)

var (
	defaultOptions = []Option{
		WithAddr("0.0.0.0:8081"),
		WithGRPCClientConfig(&config.GRPCClient{
			Addrs: []string{
				"0.0.0.0:8081",
			},
		}),
	}
)

// WithAddr returns Option that sets addr.
func WithAddr(addr string) Option {
	return func(c *gatewayClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

// WithGRPCClientConfig returns Option that sets config.
func WithGRPCClientConfig(cfg *config.GRPCClient) Option {
	return func(c *gatewayClient) {
		if cfg != nil {
			c.cfg = cfg.Bind()
		}
	}
}
