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

// Package grpc provides grpc client functions
package grpc

import (
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Option func(*ngtdClient)

var defaultOptions = []Option{
	WithAddr("127.0.0.1:8200"),
	WithGRPCClientOption(
		(&config.GRPCClient{
			Addrs: []string{
				"127.0.0.1:8200",
			},
			CallOption: &config.CallOption{
				MaxRecvMsgSize: 100000000000,
			},
			DialOption: &config.DialOption{
				Insecure: true,
			},
		}).Bind().Opts()...,
	),
}

func WithAddr(addr string) Option {
	return func(c *ngtdClient) {
		if len(addr) != 0 {
			c.addr = addr
		}
	}
}

func WithGRPCClientOption(opts ...grpc.Option) Option {
	return func(c *ngtdClient) {
		if len(opts) != 0 {
			c.opts = opts
		}
	}
}
