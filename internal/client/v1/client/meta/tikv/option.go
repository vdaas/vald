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

package tikv

import (
    "github.com/vdaas/vald/internal/net/grpc"
)

type Option func(*client) error

var defaultOptions = []Option{WithRegionErrorRetryLimit(3)}

// WithAddrs sets the TiKV store addresses (fallback RoundRobin list).
func WithAddrs(addrs ...string) Option {
    return func(c *client) error {
        if len(addrs) == 0 {
            return nil
        }
        c.addrs = append(c.addrs, addrs...)
        return nil
    }
}

// WithClient sets a prepared grpc.Client.
func WithClient(cl grpc.Client) Option {
    return func(c *client) error {
        if cl != nil {
            c.c = cl
        }
        return nil
    }
}

// WithPDClient injects existing PD client.
func WithPDClient(pc grpc.Client) Option {
    return func(c *client) error {
        if pc != nil {
            c.pd.c = pc
        }
        return nil
    }
}

// WithPDAddrs creates PD client internally from addresses.
func WithPDAddrs(addrs ...string) Option {
    return func(c *client) error {
        if len(addrs) == 0 {
            return nil
        }
        c.addrs = append(c.addrs, addrs...)
        return nil
    }
}

// WithRegionErrorRetryLimit sets the retry limit for region errors.
func WithRegionErrorRetryLimit(limit int) Option {
		return func(c *client) error {
				if limit > 0 {
						c.regionErrorRetryLimit = limit
				}
				return nil
		}
}
