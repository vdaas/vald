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

// Package ingress provides ingress filter client logic
package ingress

import "github.com/vdaas/vald/internal/net/grpc"

type Option func(*client) error

var defaultOptions = []Option{}

func WithAddrs(addrs ...string) Option {
	return func(c *client) error {
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

func WithClient(cl grpc.Client) Option {
	return func(c *client) error {
		if cl != nil {
			c.c = cl
		}
		return nil
	}
}
