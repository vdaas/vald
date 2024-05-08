//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package core provides agent ngt gRPC client functions
package core

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Option is agentClient configure.
type Option func(*agentClient) error

var defaultOptions = []Option{}

// WithAddr returns Option that sets addr.
func WithAddrs(addrs ...string) Option {
	return func(c *agentClient) error {
		if addrs == nil {
			return nil
		}
		if c.addrs != nil {
			c.addrs = append(c.addrs, addrs...)
		} else {
			c.addrs = addrs
		}
		return nil
	}
}

func WithValdClient(vc vald.Client) Option {
	return func(c *agentClient) error {
		if vc != nil {
			c.Client = vc
			if c.c != nil {
				err := c.c.Close(context.Background())
				if err != nil {
					return err
				}
			}
			c.c = c.Client.GRPCClient()
		}
		return nil
	}
}

func WithGRPCClient(g grpc.Client) Option {
	return func(c *agentClient) error {
		if g != nil {
			c.c = g
		}
		return nil
	}
}
